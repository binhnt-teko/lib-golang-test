#define _GNU_SOURCE

#include <inttypes.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static int exit_code = 0;

static void print_signals(uint64_t bitmap)
{
    for (int sig = 1; sig < NSIG; ++sig) {
        if (bitmap & (1ULL << (sig - 1))) {
            const char *sig_name = sigabbrev_np(sig);
            if (sig_name) {
                printf("%s\n", sig_name);
            } else {
                fprintf(stderr, "Unknown signal: %d\n",
                        sig);
                exit_code = 1;
            }
        }
    }
}

int main(int argc, char *argv[])
{
    uint64_t bitmap;

    if (argc != 2) {
        fprintf(stderr, "Usage: %s <bitmap>\n", argv[0]);
        return EXIT_FAILURE;
    }

    if (sscanf(argv[1], "%" SCNx64, &bitmap) != 1) {
        fprintf(stderr, "Invalid signal bitmap hex: %s\n",
                argv[1]);
        return EXIT_FAILURE;
    }

    print_signals(bitmap);

    return exit_code;
}