#include <unistd.h>
#include <sys/syscall.h> // contiene las llamadas al sistema
#include <stdio.h>

int main() {
    long id = syscall(SYS_gettid); // obtiene el id del proceso actual
    printf("Thread ID: %ld\n", id);
    return 0;
}