#ifndef PLAYER_H
#define PLAYER_H
#include <QString>

#include "bass.h"

class Player
{
public:
    Player();
    ~Player();
    bool    devInit();
    QString getTags();
    void    setVol(int vol);
    int     getVol();
    QString openAudio(const QString &uri);
    void    play();
    void    pause();
    void    stop();
    int     getPos();
    void    setPos(int pos);
    int     getBitRate();
    void    getFFT(float *array);
    bool    isPlaying();
    int     getFreq();
    void    setFreq(float freq);
    void    eqReady();
    void    disableEQ();
    void    setEQ(int id, int gain);
    QString getNowPlayInfo();
    QString getTotalTime();
    QString getCurTime();
    int     getCurTimeMS();
    int     getTotalTimeMS();
    double  getCurTimeSec();
    double  getTotalTimeSec();
    DWORD   getLevel();
    QString getFileTotalTime(const QString &fileName);
    double  getFileSecond(const QString &fileName);
    void    setReverse(bool isEnable);
    void    updateReverb(int value);
    void    setJumpPoint(double timeFrom, double timeTo);
    void    removeJumpPoint();

private:
    HSTREAM m_hNowPlay;
    HFX     m_hEqFX;
    HFX     m_hReverbFX;
    HSYNC   m_hJumping;
    bool    m_bPlayNextEnable {true};
};
#endif // PLAYER_H
