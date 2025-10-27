#include <evdi_lib.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <signal.h>
#include <sys/select.h>
#include <errno.h>

static volatile int running = 1;

// Gestionnaire de signal pour arrêt propre
void signal_handler(int signum)
{
    printf("\nArrêt demandé...\n");
    running = 0;
}

// Structure pour le contexte
typedef struct
{
    evdi_handle handle;
    void *buffer_memory;
    int buffer_id;
    int width;
    int height;
    int buffer_registered;
} app_context;

// Callback mode changé
void mode_changed_handler(struct evdi_mode mode, void *user_data)
{
    app_context *ctx = (app_context *)user_data;
    printf("Mode: %dx%d @%dHz\n", mode.width, mode.height, mode.refresh_rate);

    // Mise à jour des dimensions si nécessaire
    if (mode.width > 0 && mode.height > 0)
    {
        ctx->width = mode.width;
        ctx->height = mode.height;
    }
}

// Callback DPMS
void dpms_handler(int dpms_mode, void *user_data)
{
    const char *modes[] = {"On", "Standby", "Suspend", "Off"};
    printf("DPMS: %s\n", modes[dpms_mode]);
}

// Callback CRTC
void crtc_state_handler(int state, void *user_data)
{
    printf("CRTC: %s\n", state ? "Actif" : "Inactif");
}

// Callback update ready
void update_ready_handler(int buffer_id, void *user_data)
{
    app_context *ctx = (app_context *)user_data;

    if (buffer_id != ctx->buffer_id)
    {
        return;
    }

    // Récupérer les rectangles modifiés
    struct evdi_rect rects[16];
    int num_rects = 16;

    evdi_grab_pixels(ctx->handle, rects, &num_rects);
}

// Callback curseur (silencieux)
void cursor_set_handler(struct evdi_cursor_set cursor_set, void *user_data)
{
    if (cursor_set.buffer)
    {
        free(cursor_set.buffer);
    }
}

void cursor_move_handler(struct evdi_cursor_move cursor_move, void *user_data)
{
    // Silencieux
}

// Callback DDC/CI
void ddcci_data_handler(struct evdi_ddcci_data ddcci_data, void *user_data)
{
    // Silencieux
}

int main(int argc, char *argv[])
{
    printf("=== EVDI Écran virtuel vide (version sûre) ===\n\n");

    // Signaux
    signal(SIGINT, signal_handler);
    signal(SIGTERM, signal_handler);

    // Vérifier la version de la bibliothèque
    struct evdi_lib_version version;
    evdi_get_lib_version(&version);
    printf("Version libevdi: %d.%d.%d\n",
           version.version_major, version.version_minor, version.version_patchlevel);

    // Contexte
    app_context ctx = {0};
    ctx.width = 1920;
    ctx.height = 1080;
    ctx.buffer_id = 0;
    ctx.buffer_registered = 0;

    // Ouvrir un périphérique EVDI
    printf("Ouverture du périphérique EVDI...\n");
    ctx.handle = evdi_open_attached_to(NULL);

    if (!ctx.handle || ctx.handle == EVDI_INVALID_HANDLE)
    {
        fprintf(stderr, "ERREUR: Impossible d'ouvrir EVDI\n");
        fprintf(stderr, "- Chargez le module: sudo modprobe evdi\n");
        fprintf(stderr, "- Lancez avec sudo: sudo %s\n", argv[0]);
        return 1;
    }

    printf("✓ Périphérique ouvert\n");

    // EDID minimal mais valide pour 1920x1080@60Hz
    unsigned char edid[128];
    memset(edid, 0, sizeof(edid));

    // Header EDID
    unsigned char edid_header[] = {
        0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
        0x10, 0xAC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
        0x01, 0x1E, 0x01, 0x03, 0x80, 0x00, 0x00, 0x78,
        0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00};
    memcpy(edid, edid_header, sizeof(edid_header));

    // Timing 1920x1080@60Hz (descriptor 1)
    unsigned char timing[] = {
        0x02, 0x3A, 0x80, 0x18, 0x71, 0x38, 0x2D, 0x40,
        0x58, 0x2C, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00,
        0x00, 0x1E};
    memcpy(&edid[54], timing, sizeof(timing));

    // Checksum
    int sum = 0;
    for (int i = 0; i < 127; i++)
    {
        sum += edid[i];
    }
    edid[127] = (256 - (sum % 256)) % 256;

    // Connecter
    evdi_connect(ctx.handle, edid, sizeof(edid), 0xFFFFFFFF);
    printf("✓ Écran connecté (1920x1080@60Hz)\n");

    // Attendre un peu que le système détecte l'écran
    sleep(1);

    // Allouer buffer
    size_t buffer_size = ctx.width * ctx.height * 4;
    ctx.buffer_memory = calloc(1, buffer_size);

    if (!ctx.buffer_memory)
    {
        fprintf(stderr, "ERREUR: Allocation mémoire échouée\n");
        evdi_disconnect(ctx.handle);
        evdi_close(ctx.handle);
        return 1;
    }

    printf("✓ Buffer alloué (%zu MB)\n", buffer_size / (1024 * 1024));

    // Enregistrer buffer
    struct evdi_buffer buffer;
    buffer.id = ctx.buffer_id;
    buffer.buffer = ctx.buffer_memory;
    buffer.width = ctx.width;
    buffer.height = ctx.height;
    buffer.stride = ctx.width * 4;
    buffer.rects = NULL;
    buffer.rect_count = 0;

    evdi_register_buffer(ctx.handle, buffer);
    ctx.buffer_registered = 1;
    printf("✓ Buffer enregistré\n");

    // Configuration des événements
    struct evdi_event_context event_ctx;
    memset(&event_ctx, 0, sizeof(event_ctx));
    event_ctx.mode_changed_handler = mode_changed_handler;
    event_ctx.dpms_handler = dpms_handler;
    event_ctx.crtc_state_handler = crtc_state_handler;
    event_ctx.update_ready_handler = update_ready_handler;
    event_ctx.cursor_set_handler = cursor_set_handler;
    event_ctx.cursor_move_handler = cursor_move_handler;
    event_ctx.ddcci_data_handler = ddcci_data_handler;
    event_ctx.user_data = &ctx;

    printf("\n=== Écran virtuel actif ===\n");
    printf("État: Fond noir affiché\n");
    printf("Action: Utilisez xrandr pour configurer l'écran\n");
    printf("Arrêt: Ctrl+C\n\n");

    // Boucle événements (simple et non-bloquante)
    int idle_count = 0;

    while (running)
    {
        // Traiter les événements disponibles
        int events_handled = 0;

        // Essayer de traiter les événements
        for (int i = 0; i < 10 && running; i++)
        {
            evdi_handle_events(ctx.handle, &event_ctx);
        }

        // Pause courte pour ne pas saturer le CPU
        usleep(20000); // 20ms = 50 FPS max

        // Affichage de vie
        if (++idle_count >= 500)
        { // ~10 secondes
            printf("Actif... (Ctrl+C pour quitter)\n");
            idle_count = 0;
        }
    }

    // Nettoyage
    printf("\n=== Nettoyage ===\n");

    if (ctx.buffer_registered)
    {
        evdi_unregister_buffer(ctx.handle, ctx.buffer_id);
        printf("✓ Buffer désenregistré\n");
    }

    if (ctx.buffer_memory)
    {
        free(ctx.buffer_memory);
        printf("✓ Mémoire libérée\n");
    }

    evdi_disconnect(ctx.handle);
    printf("✓ Écran déconnecté\n");

    evdi_close(ctx.handle);
    printf("✓ Périphérique fermé\n");

    printf("\nTerminé proprement.\n");
    return 0;
}