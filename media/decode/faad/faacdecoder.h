#ifndef FAACDecoder_h
#define FAACDecoder_h

int get_one_ADTS_frame(unsigned char* buffer, size_t buf_size, unsigned char* data ,size_t* data_size);
void* faad_decoder_create(int sample_rate, int channels, int bit_rate);
int faad_decode_frame(void *pParam, unsigned char *pData, int nLen, unsigned char *pPCM, unsigned int *outLen);
void faad_decode_close(void *pParam);

#endif /* FAACDecoder_h */