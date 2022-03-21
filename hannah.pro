TEMPLATE = subdirs

SUBDIRS = \
    desktop

win32: SUBDIRS += registerProtocolHandler
    
