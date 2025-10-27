#include <evdi_lib.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <signal.h>
#include <poll.h>

static volatile int running = 1;

// Gestionnaire de signal pour arrêt propre
void signal_handler(int signum)
{
    printf("\nArrêt du programme...\n");
    running = 0;
}

// Structure pour stocker le contexte de l'application
typedef struct
{
    evdi_handle handle;
    void *buffer_memory;
    int buffer_id;
    int width;
    int height;
} app_context;

// Callback appelé quand le mode d'affichage change
void mode_changed_handler(struct evdi_mode mode, void *user_data)
{
    app_context *ctx = (app_context *)user_data;
    printf("Mode changé: %dx%d @%dHz (bpp: %d, format: 0x%x)\n",
           mode.width, mode.height, mode.refresh_rate,
           mode.bits_per_pixel, mode.pixel_format);

    // 1. Libérer l'ancien buffer
    evdi_unregister_buffer(ctx->handle, ctx->buffer_id);
    if (ctx->buffer_memory)
    {
        free(ctx->buffer_memory);
    }

    // 2. Mettre à jour les dimensions
    ctx->width = mode.width;
    ctx->height = mode.height;

    // 3. Allouer le nouveau buffer (toujours noir ici)
    ctx->buffer_memory = calloc(ctx->height * ctx->width * 4, 1);
    if (!ctx->buffer_memory)
    {
        fprintf(stderr, "Erreur: réallocation mémoire échouée\n");
        running = 0; // Arrêt propre
        return;
    }

    // 4. Réenregistrer le nouveau buffer
    struct evdi_buffer buffer;
    buffer.id = ctx->buffer_id;
    buffer.buffer = ctx->buffer_memory;
    buffer.width = ctx->width;
    buffer.height = ctx->height;
    buffer.stride = ctx->width * 4;
    buffer.rects = NULL;
    buffer.rect_count = 0;

    evdi_register_buffer(ctx->handle, buffer);
    printf("Nouveau buffer %dx%d enregistré.\n", ctx->width, ctx->height);
}

// Callback pour les événements DPMS (gestion de l'alimentation)
void dpms_handler(int dpms_mode, void *user_data)
{
    const char *mode_str[] = {"On", "Standby", "Suspend", "Off"};
    printf("DPMS: %s\n", mode_str[dpms_mode]);
}

// Callback pour les changements d'état CRTC
void crtc_state_handler(int state, void *user_data)
{
    printf("État CRTC: %s\n", state ? "Actif" : "Inactif");
}

// Callback pour les demandes de mise à jour
void update_ready_handler(int buffer_to_be_updated, void *user_data)
{
    app_context *ctx = (app_context *)user_data;

    struct evdi_rect rects[16];
    int num_rects = 16;

    // 1. Récupérer les pixels du kernel
    evdi_grab_pixels(ctx->handle, rects, &num_rects);

    if (num_rects > 0)
    {
        printf("Mise à jour reçue pour %d rectangles\n", num_rects);

        // 2. Si vous modifiiez l'image (non fait ici, reste noir), ce serait le moment.
        // Par exemple: dessiner quelque chose dans ctx->buffer_memory.

        // *************************************************************
        // 3. ENVOYER la prochaine image.
        // C'est la confirmation qui met fin au "flip" et évite le timeout.
        // *************************************************************
        struct evdi_buffer buffer;
        buffer.id = ctx->buffer_id;
        buffer.buffer = ctx->buffer_memory;
        buffer.width = ctx->width;
        buffer.height = ctx->height;
        buffer.stride = ctx->width * 4;
        buffer.rects = rects;          // Optionnel : ne notifier que les zones modifiées
        buffer.rect_count = num_rects; // (si vous avez mis à jour l'image vous-même)

        evdi_register_buffer(ctx->handle, buffer);
    }
    // Si num_rects est 0, on n'a rien à faire.
}

// Callback pour les changements de curseur
void cursor_set_handler(struct evdi_cursor_set cursor_set, void *user_data)
{
    printf("Curseur: hotspot=(%d,%d) %dx%d %s\n",
           //    cursor_set.x, cursor_set.y,
           cursor_set.hot_x, cursor_set.hot_y,
           cursor_set.width, cursor_set.height,
           cursor_set.enabled ? "visible" : "caché");

    // Libérer le buffer du curseur si présent
    if (cursor_set.buffer)
    {
        free(cursor_set.buffer);
    }
}

// Callback pour les mouvements de curseur
void cursor_move_handler(struct evdi_cursor_move cursor_move, void *user_data)
{
    // Affichage silencieux pour éviter le spam
    static int count = 0;
    if (++count % 100 == 0)
    {
        printf("Curseur déplacé (dernière pos: %d, %d)\n",
               cursor_move.x, cursor_move.y);
    }
}

// Callback pour DDC/CI
void ddcci_data_handler(struct evdi_ddcci_data ddcci_data, void *user_data)
{
    printf("Requête DDC/CI reçue: adresse=0x%02x, flags=0x%x\n",
           ddcci_data.address, ddcci_data.flags);
}

int main(int argc, char *argv[])
{
    printf("=== Programme d'exemple libevdi v1.14+ - Écran virtuel vide ===\n\n");

    // Installer le gestionnaire de signal
    signal(SIGINT, signal_handler);
    signal(SIGTERM, signal_handler);

    // Vérifier la version de la bibliothèque
    struct evdi_lib_version version;
    evdi_get_lib_version(&version);
    printf("Version libevdi: %d.%d.%d\n",
           version.version_major, version.version_minor, version.version_patchlevel);

    // Initialiser le contexte de l'application
    app_context ctx = {0};
    ctx.width = 1920;
    ctx.height = 1080;
    ctx.buffer_id = 0;

    // Méthode moderne (v1.9.0+) : ouvrir/créer automatiquement un noeud EVDI
    printf("Ouverture d'un périphérique EVDI...\n");
    ctx.handle = evdi_open_attached_to_fixed(NULL, 0);

    if (!ctx.handle || ctx.handle == EVDI_INVALID_HANDLE)
    {
        fprintf(stderr, "Erreur: impossible d'ouvrir un périphérique EVDI\n");
        fprintf(stderr, "Vérifiez que le module est chargé: sudo modprobe evdi\n");
        fprintf(stderr, "Et que vous avez les droits suffisants (sudo).\n");
        return 1;
    }

    printf("Périphérique EVDI ouvert avec succès\n");

    // // EDID pour un écran 1920x1080 @60Hz
    // // Cet EDID est un exemple simplifié mais valide
    // unsigned char edid[128] = {
    //     0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, // Header
    //     0x10, 0xAC, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Manufacturer
    //     0x01, 0x1E,                                     // Version 1.4
    //     0x80,                                           // Digital input
    //     0x00, 0x00,                                     // Screen size
    //     0x78,                                           // Features
    //     0x0A,                                           // Color characteristics
    //     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x00,
    //     // Descriptor 1: 1920x1080p60
    //     0x02, 0x3A, 0x80, 0x18, 0x71, 0x38, 0x2D, 0x40,
    //     0x58, 0x2C, 0x45, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x1E,
    //     // Descriptors 2-4: Dummy
    //     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x00,
    //     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x00,
    //     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    //     0x00, 0x00,
    //     0x00, 0x00 // Padding + checksum (à calculer)
    // };

    // // Calculer le checksum EDID (somme doit être 0 mod 256)
    // int sum = 0;
    // for (int i = 0; i < 127; i++)
    // {
    //     sum += edid[i];
    // }
    // edid[127] = (256 - (sum % 256)) % 256;
    const unsigned char edid[256] = {
        0,
        0377,
        0377,
        0377,
        0377,
        0377,
        0377,
        0,
        'N',
        0215,
        'T',
        '$',
        0,
        0,
        0,
        0,
        010,
        026,
        01,
        04,
        0245,
        'f',
        '9',
        'x',
        '"',
        015,
        0311,
        0240,
        'W',
        'G',
        0230,
        047,
        022,
        'H',
        'L',
        0277,
        0357,
        0200,
        0201,
        0200,
        0201,
        '@',
        'q',
        'O',
        0201,
        0,
        0263,
        0,
        0225,
        0,
        0225,
        017,
        0251,
        '@',
        02,
        ':',
        0200,
        030,
        'q',
        '8',
        '-',
        '@',
        'X',
        ',',
        'E',
        0,
        0372,
        '<',
        '2',
        0,
        0,
        036,
        'f',
        '!',
        'P',
        0260,
        'Q',
        0,
        033,
        '0',
        '@',
        'p',
        '6',
        0,
        0372,
        '<',
        '2',
        0,
        0,
        036,
        0,
        0,
        0,
        0375,
        0,
        '8',
        'K',
        036,
        'Q',
        021,
        0,
        012,
        ' ',
        ' ',
        ' ',
        ' ',
        ' ',
        ' ',
        0,
        0,
        0,
        0374,
        0,
        'A',
        't',
        'h',
        'e',
        'n',
        'a',
        'D',
        'P',
        012,
        ' ',
        ' ',
        ' ',
        ' ',
        01,
        '5',
        02,
        03,
        022,
        0301,
        'E',
        0220,
        037,
        04,
        023,
        03,
        '#',
        011,
        07,
        07,
        0203,
        01,
        0,
        0,
        02,
        ':',
        0200,
        0320,
        'r',
        '8',
        '-',
        '@',
        020,
        ',',
        'E',
        0200,
        0372,
        '<',
        '2',
        0,
        0,
        036,
        01,
        035,
        0,
        'r',
        'Q',
        0320,
        036,
        ' ',
        'n',
        '(',
        'U',
        0,
        0372,
        '<',
        '2',
        0,
        0,
        036,
        01,
        035,
        0,
        0274,
        'R',
        0320,
        036,
        ' ',
        0270,
        '(',
        'U',
        '@',
        0372,
        '<',
        '2',
        0,
        0,
        036,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0,
        0235,
    };

    // Connecter l'écran virtuel
    evdi_connect(ctx.handle, edid, sizeof(edid), 0xFFFFFFFF);
    printf("Écran virtuel connecté avec EDID 1920x1080@60Hz\n");

    // Allouer le buffer pour l'écran (noir par défaut)
    ctx.buffer_memory = calloc(ctx.height * ctx.width * 4, 1);
    if (!ctx.buffer_memory)
    {
        fprintf(stderr, "Erreur: allocation mémoire échouée\n");
        evdi_close(ctx.handle);
        return 1;
    }

    // Enregistrer le buffer
    struct evdi_buffer buffer;
    buffer.id = ctx.buffer_id;
    buffer.buffer = ctx.buffer_memory;
    buffer.width = ctx.width;
    buffer.height = ctx.height;
    buffer.stride = ctx.width * 4; // 4 bytes par pixel (BGRA)
    buffer.rects = NULL;
    buffer.rect_count = 0;

    evdi_register_buffer(ctx.handle, buffer);
    printf("Buffer %dx%d enregistré (écran noir)\n", ctx.width, ctx.height);

    // Configurer les callbacks
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
    printf("Appuyez sur Ctrl+C pour quitter.\n");
    printf("L'écran affiche un fond noir.\n");
    printf("Utilisez xrandr pour configurer la résolution.\n\n");

    // Obtenir le file descriptor pour les événements
    evdi_selectable event_fd = evdi_get_event_ready(ctx.handle);

    // Boucle principale d'événements avec poll
    while (running)
    {
        struct pollfd fds[1];
        fds[0].fd = event_fd;
        fds[0].events = POLLIN;

        int ret = poll(fds, 1, 100); // Timeout de 100ms

        if (ret > 0 && (fds[0].revents & POLLIN))
        {
            evdi_handle_events(ctx.handle, &event_ctx);
        }

        // Demander une mise à jour périodique
        static int frame_count = 0;
        if (++frame_count % 60 == 0)
        {
            evdi_request_update(ctx.handle, ctx.buffer_id);
        }
    }

    // Nettoyage
    printf("\n=== Nettoyage ===\n");
    evdi_unregister_buffer(ctx.handle, ctx.buffer_id);
    free(ctx.buffer_memory);

    printf("Déconnexion de l'écran virtuel...\n");
    evdi_disconnect(ctx.handle);
    evdi_close(ctx.handle);

    printf("Programme terminé proprement.\n");
    return 0;
}