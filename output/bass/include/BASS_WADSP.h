// Written by Bernd Niedergesaess
// Version: 2.4.1.0
// Copyright: 2005-2009, radio42, Bernd Niedergesaess
// All Rights reserved.
//
// Note: 
// You should disable floating-point exceptions in your app.
// This is because some Winamp DSPs might change the FloatingPointUnit state and raise a stupid exception.
// Simply call before using/loading this library: 
// "_control87(-1,_MCW_EM);"
#ifndef BASS_WADSP_H
#define BASS_WADSP_H

#include <bass.h>

#ifdef __cplusplus
extern "C" {
#endif

#ifndef BASSDSPDEF
#define BASSDSPDEF(f) WINAPI f
#endif

typedef DWORD HWADSP;

// Winamp SDK message parameter values (for lParam)
#define BASS_WADSP_IPC_GETOUTPUTTIME      105
#define BASS_WADSP_IPC_ISPLAYING          104
#define BASS_WADSP_IPC_GETVERSION         0
#define BASS_WADSP_IPC_STARTPLAY          102
#define BASS_WADSP_IPC_GETINFO            126
#define BASS_WADSP_IPC_GETLISTLENGTH      124
#define BASS_WADSP_IPC_GETLISTPOS         125
#define BASS_WADSP_IPC_GETPLAYLISTFILE    211
#define BASS_WADSP_IPC_GETPLAYLISTTITLE   212
#define BASS_WADSP_IPC                    WM_USER

int BASSDSPDEF(BASS_WADSP_GetVersion)(void);
BOOL BASSDSPDEF(BASS_WADSP_Init)(HWND hwndMain);
BOOL BASSDSPDEF(BASS_WADSP_Free)(void);
BOOL BASSDSPDEF(BASS_WADSP_FreeDSP)(HWADSP plugin);
HWND BASSDSPDEF(BASS_WADSP_GetFakeWinampWnd)(HWADSP plugin);
BOOL BASSDSPDEF(BASS_WADSP_SetSongTitle)(HWADSP plugin, const char *thetitle);
BOOL BASSDSPDEF(BASS_WADSP_SetFileName)(HWADSP plugin, const char *thefile);

typedef LRESULT (CALLBACK WINAMPWINPROC)(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam);

HWADSP BASSDSPDEF(BASS_WADSP_Load)(const char *dspfile, int x, int y, int width, int height, WINAMPWINPROC *proc);
BOOL BASSDSPDEF(BASS_WADSP_Config)(HWADSP plugin);
BOOL BASSDSPDEF(BASS_WADSP_Start)(HWADSP plugin, DWORD module, DWORD hchan);
BOOL BASSDSPDEF(BASS_WADSP_Stop)(HWADSP plugin);
BOOL BASSDSPDEF(BASS_WADSP_SetChannel)(HWADSP plugin, DWORD hchan);
DWORD BASSDSPDEF(BASS_WADSP_GetModule)(HWADSP plugin);
HDSP BASSDSPDEF(BASS_WADSP_ChannelSetDSP)(HWADSP plugin, DWORD hchan, int priority);
BOOL BASSDSPDEF(BASS_WADSP_ChannelRemoveDSP)(HWADSP plugin);

DWORD BASSDSPDEF(BASS_WADSP_ModifySamplesSTREAM)(HWADSP plugin, void *buffer, DWORD length);
DWORD BASSDSPDEF(BASS_WADSP_ModifySamplesDSP)(HWADSP plugin, void *buffer, DWORD length);

LPTSTR BASSDSPDEF(BASS_WADSP_GetName)(HWADSP plugin);
UINT BASSDSPDEF(BASS_WADSP_GetModuleCount)(HWADSP plugin);
LPTSTR BASSDSPDEF(BASS_WADSP_GetModuleName)(HWADSP plugin, DWORD module);

BOOL BASSDSPDEF(BASS_WADSP_PluginInfoFree)(void);
BOOL BASSDSPDEF(BASS_WADSP_PluginInfoLoad)(const char *dspfile);
LPTSTR BASSDSPDEF(BASS_WADSP_PluginInfoGetName)(void);
UINT BASSDSPDEF(BASS_WADSP_PluginInfoGetModuleCount)(void);
LPTSTR BASSDSPDEF(BASS_WADSP_PluginInfoGetModuleName)(DWORD module);


#ifdef __cplusplus
}
#endif

#endif