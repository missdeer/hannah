#include "_cgo_export.h"

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include "neaacdec.h"

typedef struct {
    NeAACDecHandle handle;
    int sample_rate;
    int channels;
    int bit_rate;
}FAADContext;

uint32_t _get_frame_length(const unsigned char *aac_header)
{
    uint32_t len = *(uint32_t *)(aac_header + 3);
    len = ntohl(len); //Little Endian
    len = len << 6;
    len = len >> 19;
    return len;
}

void* faad_decoder_create(int sample_rate, int channels, int bit_rate)
{
    NeAACDecHandle handle = NeAACDecOpen();
    if(!handle){
        printf("NeAACDecOpen failed\n");
        goto error;
    }
    NeAACDecConfigurationPtr conf = NeAACDecGetCurrentConfiguration(handle);
    if(!conf){
        printf("NeAACDecGetCurrentConfiguration failed\n");
        goto error;
    }
    conf->defSampleRate = sample_rate;
    conf->outputFormat = FAAD_FMT_16BIT;
    conf->dontUpSampleImplicitSBR = 1;
    NeAACDecSetConfiguration(handle, conf);

    FAADContext* ctx = malloc(sizeof(FAADContext));
    ctx->handle = handle;
    ctx->sample_rate = sample_rate;
    ctx->channels = channels;
    ctx->bit_rate = bit_rate;
    return ctx;

error:
    if(handle){
        NeAACDecClose(handle);
    }
    return NULL;
}

int faad_decode_frame(void *pParam, unsigned char *pData, int nLen, unsigned char *pPCM, unsigned int *outLen)
{
    FAADContext* pCtx = (FAADContext*)pParam;
    NeAACDecHandle handle = pCtx->handle;
    long res = NeAACDecInit(handle, pData, nLen, (unsigned long*)&pCtx->sample_rate, (unsigned char*)&pCtx->channels);
    if (res < 0) {
        printf("NeAACDecInit failed\n");
        return -1;
    }
    NeAACDecFrameInfo info;
    uint32_t framelen = _get_frame_length(pData);
    unsigned char *buf = (unsigned char *)NeAACDecDecode(handle, &info, pData, framelen);
    if (buf && info.error == 0) {
        if (info.samplerate == 44100) {
            //src: 2048 samples, 4096 bytes
            //dst: 2048 samples, 4096 bytes
            int tmplen = (int)info.samples * 16 / 8;
            memcpy(pPCM,buf,tmplen);
            *outLen = tmplen;
        } else if (info.samplerate == 22050) {
            //src: 1024 samples, 2048 bytes
            //dst: 2048 samples, 4096 bytes
            short *ori = (short*)buf;
            short tmpbuf[info.samples * 2];
            int tmplen = (int)info.samples * 16 / 8 * 2;
            for (int32_t i = 0, j = 0; i < info.samples; i += 2) {
                tmpbuf[j++] = ori[i];
                tmpbuf[j++] = ori[i + 1];
                tmpbuf[j++] = ori[i];
                tmpbuf[j++] = ori[i + 1];
            }
            memcpy(pPCM,tmpbuf,tmplen);
            *outLen = tmplen;
        }else if(info.samplerate == 8000){
            //从双声道的数据中提取单通道
            for(int i=0,j=0; i<4096 && j<2048; i+=4, j+=2)
            {
                pPCM[j]= buf[i];
                pPCM[j+1]=buf[i+1];
            }
            *outLen = (unsigned int)info.samples;
        }
    } else {
        printf("NeAACDecDecode failed\n");
        return -1;
    }
    return 0;
}

void faad_decode_close(void *pParam)
{
    if(!pParam){
        return;
    }
    FAADContext* pCtx = (FAADContext*)pParam;
    if(pCtx->handle){
        NeAACDecClose(pCtx->handle);
    }
    free(pCtx);
}