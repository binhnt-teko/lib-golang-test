#include <assert.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>
#include <unistd.h>

#define EXIT_STRING "EXIT\n"

__attribute__((noinline)) static void run_child(void)
{
    char input[sizeof(EXIT_STRING)];

    while (1) {
        if (fgets(input, sizeof(input), stdin) == NULL) {
            if (ferror(stdin)) {
                fprintf(stderr,
                        "Error reading stdin in child\n");
            }
            if (feof(stdin)) {
                fprintf(stderr,
                        "EOF waiting for exit string\n");
            }
            break;
        }
        if (strcmp(input, EXIT_STRING) == 0) {
            break;
        }
        if (strlen(input) == sizeof(input) - 1 &&
            input[strlen(input) - 1] != '\n') {
            /* Partial line read, clear it */
            int c;
            while ((c = getchar()) != '\n' && c != EOF)
                ;
        }
    }

    _exit(0);
}

int main(void)
{
    pid_t pid;
    sigset_t set;

    sigfillset(&set);
    assert(sigprocmask(SIG_BLOCK, &set, NULL) == 0);

    pid = vfork();

    if (pid == 0) {
        run_child();
    } else if (pid < 0) {
        return 1;
    }

    return 0;
}