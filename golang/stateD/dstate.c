#include <unistd.h>

__attribute__((noinline)) static void run_child(void)
{
    pause();
    _exit(0);
}

int main(void)
{
    pid_t pid = vfork();

    if (pid == 0) {
        run_child();
    } else if (pid < 0) {
        return 1;
    }

    return 0;
}