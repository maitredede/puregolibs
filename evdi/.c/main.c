#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <evdi_lib.h>
#include <stdarg.h>
#include <poll.h>
#include <time.h>

#define NUM_BUFFERS 2

// Structure pour stocker les données de contexte
typedef struct
{
    evdi_handle handle;
    struct evdi_buffer buffers[NUM_BUFFERS];
    int current_width;
    int current_height;
    int buffers_registered;
} app_context_t;

/*
 * Lit un fichier binaire (comme un EDID) et retourne un buffer alloué
 * avec son contenu. La taille est retournée via le pointeur edid_size.
 * L'appelant doit libérer le buffer retourné avec free().
 */
unsigned char *read_edid_file(const char *filename, size_t *edid_size)
{
    FILE *f = fopen(filename, "rb"); // "rb" = Read Binary
    if (!f)
    {
        perror("Erreur: Impossible d'ouvrir le fichier EDID");
        return NULL;
    }

    // Obtenir la taille du fichier
    fseek(f, 0, SEEK_END);
    *edid_size = ftell(f);
    fseek(f, 0, SEEK_SET);

    if (*edid_size == 0)
    {
        fprintf(stderr, "Erreur: Fichier EDID '%s' est vide\n", filename);
        fclose(f);
        return NULL;
    }

    // Allouer la mémoire
    unsigned char *buffer = (unsigned char *)malloc(*edid_size);
    if (!buffer)
    {
        fprintf(stderr, "Erreur: Allocation mémoire pour EDID (%zu octets)\n", *edid_size);
        fclose(f);
        return NULL;
    }

    // Lire le fichier dans le buffer
    if (fread(buffer, 1, *edid_size, f) != *edid_size)
    {
        fprintf(stderr, "Erreur: Lecture du fichier EDID '%s'\n", filename);
        fclose(f);
        free(buffer);
        return NULL;
    }

    fclose(f);
    printf("Fichier EDID '%s' lu avec succès (%zu octets)\n", filename, *edid_size);
    return buffer;
}

// Callbacks EVDI
static void mode_changed_handler(struct evdi_mode mode, void *user_data)
{
    app_context_t *ctx = (app_context_t *)user_data;
    printf("Mode changé: %dx%d @%dHz\n", mode.width, mode.height, mode.refresh_rate);

    // Désenregistrer les anciens buffers si existants
    if (ctx->buffers_registered)
    {
        for (int i = 0; i < NUM_BUFFERS; i++)
        {
            evdi_unregister_buffer(ctx->handle, i);
            if (ctx->buffers[i].buffer)
            {
                free(ctx->buffers[i].buffer);
                ctx->buffers[i].buffer = NULL;
            }
        }
        ctx->buffers_registered = 0;
    }

    // Mettre à jour les dimensions
    ctx->current_width = mode.width;
    ctx->current_height = mode.height;

    // Créer et enregistrer de nouveaux buffers
    int buffer_size = mode.width * mode.height * 4; // 4 bytes par pixel (RGBA)

    for (int i = 0; i < NUM_BUFFERS; i++)
    {
        ctx->buffers[i].id = i;
        ctx->buffers[i].buffer = malloc(buffer_size);
        ctx->buffers[i].width = mode.width;
        ctx->buffers[i].height = mode.height;
        ctx->buffers[i].stride = mode.width * 4;

        if (!ctx->buffers[i].buffer)
        {
            fprintf(stderr, "Erreur: allocation mémoire buffer %d\n", i);
            continue;
        }

        memset(ctx->buffers[i].buffer, 0, buffer_size);
        evdi_register_buffer(ctx->handle, ctx->buffers[i]);
        printf("Buffer %d enregistré (%dx%d, %d bytes)\n", i, mode.width, mode.height, buffer_size);
    }

    ctx->buffers_registered = 1;

    // **AJOUT IMPORTANT**
    // Demander la première mise à jour pour démarrer la chaîne de "flip"
    printf("Demande de la première mise à jour (buffer 0)\n");
    evdi_request_update(ctx->handle, 0);
}

static void update_ready_handler(int buffer_to_be_updated, void *user_data)
{
    app_context_t *ctx = (app_context_t *)user_data;
    printf("Buffer %d prêt pour mise à jour\n", buffer_to_be_updated);

    if (!ctx->buffers_registered || buffer_to_be_updated >= NUM_BUFFERS)
    {
        // ERREUR : Si on quitte ici sans rien faire, le timeout se produira.
        // Il faut quand même signaler qu'on a "traité" (ignoré) ce buffer.
        if (ctx->handle)
        {
            evdi_request_update(ctx->handle, buffer_to_be_updated);
        }
        return;
    }

    // Demander à EVDI de remplir le buffer avec les pixels
    int num_rects = 0;
    evdi_grab_pixels(ctx->handle, NULL, &num_rects);

    // Analyser quelques pixels pour exemple
    unsigned char *pixels = (unsigned char *)ctx->buffers[buffer_to_be_updated].buffer;

    // Ajout de vérifications pour éviter les crashs si le mode n'est pas encore défini
    if (pixels && ctx->current_height > 0 && ctx->current_width > 0)
    {
        // Récupérer le pixel au centre de l'écran
        int center_x = ctx->current_width / 2;
        int center_y = ctx->current_height / 2;
        int offset = (center_y * ctx->buffers[buffer_to_be_updated].stride) + (center_x * 4);

        // Vérification simple des limites pour éviter un segfault
        if (offset < (ctx->current_width * ctx->current_height * 4) - 4)
        {
            unsigned char b = pixels[offset];
            unsigned char g = pixels[offset + 1];
            unsigned char r = pixels[offset + 2];
            unsigned char a = pixels[offset + 3];

            printf("  → Pixel central (%d,%d): RGBA(%d, %d, %d, %d)\n",
                   center_x, center_y, r, g, b, a);
        }
    }

    //
    // ** LA CORRECTION CRUCIALE EST ICI **
    //
    // Signaler au noyau que nous avons terminé avec ce buffer
    // et que nous sommes prêts pour la prochaine mise à jour.
    evdi_request_update(ctx->handle, buffer_to_be_updated);
}

static void crtc_state_handler(int state, void *user_data)
{
    printf("État CRTC: %d\n", state);
}

static void cursor_set_handler(struct evdi_cursor_set cursor_set, void *user_data)
{
    printf("Curseur défini: %dx%d\n", cursor_set.width, cursor_set.height);
    if (cursor_set.buffer)
    {
        free(cursor_set.buffer);
    }
}

static void cursor_move_handler(struct evdi_cursor_move cursor_move, void *user_data)
{
    printf("Curseur déplacé: (%d, %d)\n", cursor_move.x, cursor_move.y);
}

static void dpms_handler(int dpms_mode, void *user_data)
{
    const char *mode_str[] = {"On", "Standby", "Suspend", "Off"};
    printf("Mode DPMS: %d (%s)\n", dpms_mode, mode_str[dpms_mode]);
}

static int get_available_evdi_card()
{
    for (int i = 0; i < 20; i++)
    {
        if (evdi_check_device(i) == AVAILABLE)
        {
            return i;
        }
    }
    evdi_add_device();
    for (int i = 0; i < 20; i++)
    {
        if (evdi_check_device(i) == AVAILABLE)
        {
            return i;
        }
    }
    return -1;
}

static void evdi_log_callback(void *user_data, const char *fmt, ...)
{
    va_list args;

    // Début de la liste d'arguments variables
    va_start(args, fmt);

    // Afficher le préfixe sur stderr
    fprintf(stderr, "[EVDI LOG] ");

    // Utiliser vfprintf pour gérer la chaîne de format et les arguments
    vfprintf(stderr, fmt, args);

    // Fin de la liste d'arguments variables
    va_end(args);

    // Assurer que le message se termine par un saut de ligne si ce n'est pas déjà le cas
    if (fmt[strlen(fmt) - 1] != '\n')
    {
        fprintf(stderr, "\n");
    }
}

int main(int argc, char *argv[])
{
    struct evdi_logging log =
        {
            evdi_log_callback,
            NULL,
        };
    evdi_set_logging(log);

    evdi_handle handle;
    struct evdi_event_context event_context;
    // int device_count;
    app_context_t ctx;

    // Initialiser le contexte
    memset(&ctx, 0, sizeof(app_context_t));
    ctx.current_width = 1920;
    ctx.current_height = 1080;

    printf("=== Programme EVDI - Écran factice %dx%d ===\n\n", ctx.current_width, ctx.current_height);

    // // Vérifier les périphériques EVDI disponibles
    // device_count = evdi_check_device(0);
    // printf("Nombre de périphériques EVDI: %d\n", device_count);

    // if (device_count < 1)
    // {
    //     printf("Ajout d'un nouveau périphérique EVDI...\n");
    //     evdi_add_device();
    //     sleep(1);
    // }
    int device = get_available_evdi_card();

    // Ouvrir le périphérique EVDI
    printf("Ouverture du périphérique EVDI %d...\n", device);
    handle = evdi_open(device);

    if (handle == EVDI_INVALID_HANDLE)
    {
        fprintf(stderr, "Erreur: Impossible d'ouvrir le périphérique EVDI\n");
        return 1;
    }

    printf("Périphérique ouvert avec succès!\n");

    ctx.handle = handle;

    // Connexion du périphérique
    printf("Connexion du périphérique...\n");
    // // EDID 128-octets standard pour 1920x1080 @ 60Hz
    // unsigned char edid[128] = {
    //     0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x1A, 0xEC, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
    //     0x00, 0x1A, 0x01, 0x03, 0x80, 0x00, 0x00, 0x78, 0x0A, 0xEE, 0x91, 0xA3, 0x54, 0x4C, 0x99, 0x26,
    //     0x0F, 0x50, 0x54, 0x21, 0x08, 0x00, 0x81, 0x80, 0x81, 0x40, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
    //     0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x02, 0x3A, 0x80, 0x18, 0x71, 0x38, 0x2D, 0x40, 0x58, 0x2C,
    //     0x45, 0x00, 0x56, 0x50, 0x51, 0x00, 0x00, 0x1E, 0x01, 0x1D, 0x00, 0x72, 0x51, 0xD0, 0x1E, 0x20,
    //     0x6E, 0x28, 0x55, 0x00, 0x56, 0x50, 0x51, 0x00, 0x00, 0x1E, 0x00, 0x00, 0x00, 0xFD, 0x00, 0x38,
    //     0x4C, 0x1E, 0x53, 0x11, 0x00, 0x0A, 0x20, 0x20, 20, 0x20, 0x20, 0x20, 0x00, 0x00, 0x00, 0xFC,
    //     0x00, 0x46, 0x61, 0x6B, 0x65, 0x4D, 0x6F, 0x6E, 0x69, 0x74, 0x6F, 0x72, 0x0A, 0x00, 0x00, 0xA9};
    // evdi_connect(handle, edid, sizeof(edid), 0);
    // Charger l'EDID depuis le fichier
    printf("Chargement de l'EDID depuis le fichier...\n");
    size_t edid_size = 0;
    // Utilise le nom de fichier que vous avez fourni
    unsigned char *edid_data = read_edid_file("EDIDv2_1920x1080", &edid_size);

    if (edid_data == NULL)
    {
        // L'erreur est déjà affichée par read_edid_file
        evdi_close(handle);
        return 1;
    }

    // Connexion du périphérique avec l'EDID chargé
    printf("Connexion du périphérique...\n");
    evdi_connect(handle, edid_data, edid_size, 0);

    // Configuration des callbacks
    memset(&event_context, 0, sizeof(event_context));
    event_context.mode_changed_handler = mode_changed_handler;
    event_context.update_ready_handler = update_ready_handler;
    event_context.crtc_state_handler = crtc_state_handler;
    event_context.cursor_set_handler = cursor_set_handler;
    event_context.cursor_move_handler = cursor_move_handler;
    event_context.dpms_handler = dpms_handler;
    event_context.user_data = &ctx;

    // Obtenir le file descriptor (FD) pour les événements EVDI
    int event_fd = evdi_get_event_ready(handle);
    if (event_fd < 0)
    {
        fprintf(stderr, "Erreur: Impossible d'obtenir le FD des événements\n");
        // ... (ajoutez votre code de nettoyage ici : free(edid_data), evdi_close, etc.)
        free(edid_data);
        evdi_disconnect(handle);
        evdi_close(handle);
        return 1;
    }

    printf("\nConfiguration du mode 1920x1080 @60Hz...\n");

    printf("\nÉcran factice actif. Écoute des événements pendant 60 secondes...\n");
    printf("(Appuyez sur Ctrl+C pour arrêter)\n\n");

    // // Gérer un premier lot d'événements (pour déclencher le mode_changed_handler initial)
    // evdi_handle_events(handle, &event_context);

    // Boucle principale basée sur poll()
    struct pollfd fds[1];
    fds[0].fd = event_fd;
    fds[0].events = POLLIN; // Attendre les événements "prêts à lire"

    printf("Démarrage de la boucle principale\n");
    time_t start_time = time(NULL);
    while (time(NULL) - start_time < 60) // Boucle pendant 60 secondes
    {
        // Attendre un événement, avec un timeout de 1 seconde (1000 ms)
        int ret = poll(fds, 1, 1000);

        if (ret < 0)
        {
            perror("Erreur poll()");
            break;
        }
        else if (ret == 0)
        {
            // Timeout - 1 seconde s'est écoulée sans événement
            printf("Toujours en écoute... (Timeout poll(), %ld/60s)\n", time(NULL) - start_time);
            continue;
        }

        // Un événement est prêt!
        if (fds[0].revents & POLLIN)
        {
            printf("Événement reçu!\n"); // Décommentez pour déboguer
            evdi_handle_events(handle, &event_context);
        }
    }

    // Nettoyage
    printf("\nDéconnexion et fermeture...\n");

    // Libérer les buffers
    if (ctx.buffers_registered)
    {
        for (int i = 0; i < NUM_BUFFERS; i++)
        {
            evdi_unregister_buffer(handle, i);
            if (ctx.buffers[i].buffer)
            {
                free(ctx.buffers[i].buffer);
            }
        }
    }

    evdi_disconnect(handle);
    evdi_close(handle);

    printf("Programme terminé.\n");
    return 0;
}