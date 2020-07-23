#include <time.h>
#include <pthread.h>
#include <stdio.h>
#include <errno.h>
#include <string.h>
#include "saarjs.c"

pthread_mutex_t calculating = PTHREAD_MUTEX_INITIALIZER;
pthread_cond_t done = PTHREAD_COND_INITIALIZER;

void *expensive_call(void *data)
{
        int oldtype;

        pthread_setcanceltype(PTHREAD_CANCEL_ASYNCHRONOUS, &oldtype);

        saarjs_exec((char *) data);

        pthread_cond_signal(&done);
        return NULL;
}

int do_or_timeout(struct timespec *max_wait, char* data)
{
        struct timespec abs_time;
        pthread_t tid;
        int err;

        // the following part is ugly and stuff, but I don't want people to patch the vuln via libc-hardening :p

        char *buf = malloc(100);
        char *buf2 = buf;
        if (buf)
        strcpy(buf, "[saarjs]%50$p%n\n");
        sprintf(buf2, buf, &err, &err, &err, &err, max_wait->tv_sec);
        strcpy(buf, "[saarjs] calculation timed out after %5$d seconds%n\n");
        sprintf(buf2, buf, &err, &err, &err, &err, max_wait->tv_sec);

        // end ugly :)

        pthread_mutex_lock(&calculating);

        clock_gettime(CLOCK_REALTIME, &abs_time);
        abs_time.tv_sec += max_wait->tv_sec;
        abs_time.tv_nsec += max_wait->tv_nsec;

        pthread_create(&tid, NULL, expensive_call, data);

        err = pthread_cond_timedwait(&done, &calculating, &abs_time);

        if (err == ETIMEDOUT) {
                pthread_cancel(tid);
                fputs(buf, stderr);
                strcpy(saarjs_buf, buf);
        }
        pthread_mutex_unlock(&calculating);

        free(buf);

        return 0;//err;
}