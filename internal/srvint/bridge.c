#include "bridge.h"

#include <errno.h>

#include "srvint.h"

static int empty_application_callback(
    srvint_t *ctx, const uint8_t *request, size_t request_length,
    uint8_t *response, size_t response_capacity, size_t *response_length,
    void *user_data) {
  (void)ctx;
  (void)request;
  (void)request_length;
  (void)user_data;
  if (response == NULL || response_length == NULL || response_capacity == 0) {
    errno = EINVAL;
    return -1;
  }
  response[0] = 0;
  *response_length = 1;
  return 0;
}

srvint_t *mbrd_srvint_new(const char *device) {
  srvint_t *ctx = srvint_serial_new(device, 115200, 'N', 8, 1);
  if (ctx == NULL) {
    return NULL;
  }
  if (srvint_set_slave(ctx, SRVINT_DEVICE_ADDRESS) != 0 ||
      srvint_set_response_timeout(ctx, 0, 100000) != 0) {
    srvint_free(ctx);
    return NULL;
  }
  return ctx;
}

int mbrd_srvint_connect(srvint_t *ctx) { return srvint_connect(ctx); }

int mbrd_srvint_receive(srvint_t *ctx, uint8_t *request, size_t capacity) {
  return srvint_receive(ctx, request, capacity);
}

int mbrd_srvint_reply(srvint_t *ctx, const uint8_t *request, size_t length) {
  return srvint_reply(ctx, request, length, empty_application_callback, NULL);
}

int mbrd_srvint_close(srvint_t *ctx) { return srvint_close(ctx); }

void mbrd_srvint_free(srvint_t *ctx) { srvint_free(ctx); }

const char *mbrd_srvint_strerror(int error_number) {
  return srvint_strerror(error_number);
}

int mbrd_srvint_errno(void) { return errno; }
