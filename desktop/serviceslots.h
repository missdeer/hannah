#ifndef _SERVICESLOTS_H_
#define _SERVICESLOTS_H_

#include <QStringList>

void serviceSearch(const QString &s);
void serviceOpenUrl(const QString &s);
void serviceOpenLink(const QString &s);
void serviceAppendToPlaylist(const QStringList &s);
void serviceClearAndAddToPlaylist(const QStringList &s);
void serviceAppendToPlaylistFile(const QStringList &s);
void serviceClearAndAddToPlaylistFile(const QStringList &s);

#endif
