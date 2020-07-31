#ifndef BASSWINAMP_H
#define BASSWINAMP_H

#include <bass.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef BASSWINAMPDEF
#define BASSWINAMPDEF(f) WINAPI f
#endif


#define BASS_CTYPE_STREAM_WINAMP	0x10400

#define BASS_WINAMP_SYNC_BITRATE	100

// BASS_SetConfig, BASS_GetConfig flags
#define BASS_CONFIG_WINAMP_INPUT_TIMEOUT   0x10800  // Set the time to wait until timing out because
                                                   // the plugin is not using the output system

// BASS_WINAMP_FindPlugin flags
#define BASS_WINAMP_FIND_INPUT		1
#define BASS_WINAMP_FIND_RECURSIVE	4
// return value type
#define BASS_WINAMP_FIND_COMMALIST	8
// Delphi's comma list style (item1,item2,"item with , commas",item4,"item with space")
// the list ends with single NULL character

char* BASSWINAMPDEF(BASS_WINAMP_FindPlugins)(char *pluginpath, DWORD flags);
DWORD BASSWINAMPDEF(BASS_WINAMP_LoadPlugin)(char *f);
void BASSWINAMPDEF(BASS_WINAMP_UnloadPlugin)(DWORD handle);
char* BASSWINAMPDEF(BASS_WINAMP_GetName)(DWORD handle);
int BASSWINAMPDEF(BASS_WINAMP_GetVersion)(DWORD handle);
BOOL BASSWINAMPDEF(BASS_WINAMP_GetIsSeekable)(DWORD handle);
BOOL BASSWINAMPDEF(BASS_WINAMP_GetUsesOutput)(DWORD handle);
char* BASSWINAMPDEF(BASS_WINAMP_GetExtentions)(DWORD handle);
BOOL BASSWINAMPDEF(BASS_WINAMP_GetFileInfo)(char *f, char *title, int *lenms);
BOOL BASSWINAMPDEF(BASS_WINAMP_InfoDlg)(char *f, DWORD win);
void BASSWINAMPDEF(BASS_WINAMP_ConfigPlugin)(DWORD handle, DWORD win);
void BASSWINAMPDEF(BASS_WINAMP_AboutPlugin)(DWORD handle, DWORD win);
HSTREAM BASSWINAMPDEF(BASS_WINAMP_StreamCreate)(char *f, DWORD flags);

#ifdef __cplusplus
}
#endif

#endif
