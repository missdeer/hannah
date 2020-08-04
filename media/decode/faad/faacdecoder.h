#ifndef FAACDecoder_h
#define FAACDecoder_h

void* faad_decoder_create(int sample_rate, int channels, int bit_rate);
int faad_decode_frame(void *pParam, unsigned char *pData, int nLen, unsigned char *pPCM, unsigned int *outLen);
void faad_decode_close(void *pParam);

#endif /* FAACDecoder_h */