#ifndef MBRD_EMU_SRVINT_BRIDGE_H
#define MBRD_EMU_SRVINT_BRIDGE_H

#include <stddef.h>
#include <stdint.h>

typedef struct _srvint srvint_t;

srvint_t *mbrd_srvint_new(const char *device);
int mbrd_srvint_connect(srvint_t *ctx);
int mbrd_srvint_receive(srvint_t *ctx, uint8_t *request, size_t capacity);
int mbrd_srvint_reply(srvint_t *ctx, const uint8_t *request, size_t length);
int mbrd_srvint_close(srvint_t *ctx);
void mbrd_srvint_free(srvint_t *ctx);
const char *mbrd_srvint_strerror(int error_number);
int mbrd_srvint_errno(void);

#endif
