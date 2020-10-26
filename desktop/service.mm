#import <QUrl>
#import <QtCore>

#import "service.h"

#import <Cocoa/Cocoa.h>

@implementation HannahService

- (void)search:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![pboard canReadObjectForClasses:@[ [NSString class] ] options:@{}])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain a string.", nil);
        return;
    }

    // Try to parse a selector from the pasteboard contents.
    NSString *pasteboardString = [pboard stringForType:NSPasteboardTypeString];

    QString s = QString::fromNSString(pasteboardString);
    qDebug() << "search keyword:" << s;

    //    NSString *methodName = [AKMethodNameExtractor extractMethodNameFromString:pasteboardString];

    //    if (methodName == nil)
    //    {
    //        NSBeep();
    //        return;
    //    }

    //    // Stuff the extracted method name into the system paste buffer.
    //    NSPasteboard *generalPasteboard = [NSPasteboard generalPasteboard];

    //    [generalPasteboard declareTypes:@[ NSStringPboardType ] owner:nil];
    //    [generalPasteboard setString:methodName forType:NSStringPboardType];
}

- (void)openUrl:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![pboard canReadObjectForClasses:@[ [NSString class] ] options:@{}])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain a string.", nil);
        return;
    }

    // Try to parse a selector from the pasteboard contents.
    NSString *pasteboardString = [pboard stringForType:NSPasteboardTypeString];

    QString s = QString::fromNSString(pasteboardString);
    qDebug() << "open Url:" << s;
}

- (void)openLink:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![pboard canReadObjectForClasses:@[ [NSString class] ] options:@{}])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain a string.", nil);
        return;
    }

    // Try to parse a selector from the pasteboard contents.
    NSString *pasteboardString = [pboard stringForType:NSPasteboardTypeString];

    QString s = QString::fromNSString(pasteboardString);
    qDebug() << "open link:" << s;
}

- (void)appendToPlaylist:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![[pboard types] containsObject:NSFilenamesPboardType])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain directory(s)/file(s).", nil);
        return;
    }

    NSArray *fileArray = [pboard propertyListForType:NSFilenamesPboardType];
    for (NSString *filePath in fileArray)
    {
        QString s = QString::fromNSString(filePath);
        qDebug() << "appendToPlaylist:" << s;
    }
}

- (void)clearAndAddToPlaylist:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![[pboard types] containsObject:NSFilenamesPboardType])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain directory(s)/file(s).", nil);
        return;
    }

    NSArray *fileArray = [pboard propertyListForType:NSFilenamesPboardType];
    for (NSString *filePath in fileArray)
    {
        QString s = QString::fromNSString(filePath);
        qDebug() << "clearAndAddToPlaylist:" << s;
    }
}

- (void)appendToPlaylistFile:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![[pboard types] containsObject:NSFilenamesPboardType])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain directory(s)/file(s).", nil);
        return;
    }

    NSArray *fileArray = [pboard propertyListForType:NSFilenamesPboardType];
    for (NSString *filePath in fileArray)
    {
        QString s = QString::fromNSString(filePath);
        qDebug() << "appendToPlaylistFile:" << s;
    }
}

- (void)clearAndAddToPlaylistFile:(NSPasteboard *)pboard userData:(NSString *)userData error:(NSString **)error
{
    if (![[pboard types] containsObject:NSFilenamesPboardType])
    {
        *error = NSLocalizedString(@"Error: the pasteboard doesn't contain directory(s)/file(s).", nil);
        return;
    }

    NSArray *fileArray = [pboard propertyListForType:NSFilenamesPboardType];
    for (NSString *filePath in fileArray)
    {
        QString s = QString::fromNSString(filePath);
        qDebug() << "clearAndAddToPlaylistFile:" << s;
    }
}

@end

void registerHannahService()
{
    NSRegisterServicesProvider([HannahService new], @"Hannah");
    NSUpdateDynamicServices();
}
