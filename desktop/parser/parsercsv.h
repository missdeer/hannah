//
// C++ Interface: parserm3u
//
// Description: Interface header parse Comma-Separated Values (CSV) formatted playlists (rfc4180)
//
//
// Author: Ingo Kossyk <kossyki@cs.tu-berlin.de>, (C) 2004
// Author: Tobias Rafreider trafreider@mixxx.org, (C) 2011
// Author: Daniel Schürmann daschuer@gmx.de, (C) 2011
//
// Copyright: See COPYING file that comes with this distribution
//
//
#pragma once

#include <QByteArray>
#include <QList>
#include <QString>

#include "parser.h"

class ParserCsv : public Parser
{
    Q_OBJECT
public:
    ParserCsv()  = default;
    ~ParserCsv() = default;
    /**Overwriting function parse in class Parser**/
    QList<QString> parse(const QString &);

private:
    /**Reads a line from the file and returns filepath if a valid file**/
    QList<QList<QString>> tokenize(const QByteArray &str, char delimiter);
};
