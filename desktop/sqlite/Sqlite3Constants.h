#pragma once

const char *const SQL_STATEMENT_CIPHER_VERSION          = "PRAGMA cipher_version;";
const char *const SQL_STATEMENT_CIPHER_PROVIDER         = "PRAGMA cipher_provider;";
const char *const SQL_STATEMENT_CIPHER_PROVIDER_VERSION = "PRAGMA cipher_provider_version;";
const char *const SQL_STATEMENT_ENABLE_FOREIGN_KEYS     = "PRAGMA foreign_keys = ON;";

// SQLite options for performance issue
const char *const SQL_STATEMENT_MMAP_SIZE    = "PRAGMA mmap_size=268435456;";
const char *const SQL_STATEMENT_SYNCHRONOUS  = "PRAGMA synchronous = OFF;";
const char *const SQL_STATEMENT_JOURNAL_MODE = "PRAGMA journal_mode = MEMORY;";
const char *const SQL_STATEMENT_CACHE_SIZE   = "PRAGMA cache_size=-400;";

// SQLCipher  options for performance issue
const char *const SQL_STATEMENT_KDF_ITER_V3_MOBILE          = "PRAGMA kdf_iter = 10000;";
const char *const SQL_STATEMENT_CIPHER_PAGE_SIZE_V3_MOBILE  = "PRAGMA cipher_page_size = 4096;";
const char *const SQL_STATEMENT_KDF_ITER_V3_DESKTOP         = "PRAGMA kdf_iter = 64000;";
const char *const SQL_STATEMENT_CIPHER_PAGE_SIZE_V3_DESKTOP = "PRAGMA cipher_page_size = 1024;";

// SQLCipher options for v3.x
const char *const SQL_STATEMENT_CIPHER_HMAC_ALGO = "PRAGMA cipher_hmac_algorithm = HMAC_SHA1;";
const char *const SQL_STATEMENT_CIPHER_KDF_ALGO  = "PRAGMA cipher_kdf_algorithm = PBKDF2_HMAC_SHA1;";

// SQLCipher options for performance issue upgrade
const char *const SQL_STATEMENT_NEWDB_KDF_ITER         = "PRAGMA newdb.kdf_iter = '10000';";
const char *const SQL_STATEMENT_NEWDB_CIPHER_PAGE_SIZE = "PRAGMA newdb.cipher_page_size = 4096;";

// export SQLCipher database
const char *const SQL_STATEMENT_ATTACH_DATABASE_NEWDB_KEY = "ATTACH DATABASE ? AS newdb KEY ?;";
const char *const SQL_STATEMENT_NEWDB_EXPORT              = "SELECT sqlcipher_export('newdb');";
const char *const SQL_STATEMENT_NEWDB_DETACH              = "DETACH DATABASE newdb;";

// for error recovery
const char *const SQL_STATEMENT_INTEGRITY_CHECK = "PRAGMA integrity_check;";

// transaction
const char *const SQL_STATEMENT_BEGIN_TRANSACTION    = "BEGIN TRANSACTION;";
const char *const SQL_STATEMENT_COMMIT_TRANSACTION   = "COMMIT TRANSACTION;";
const char *const SQL_STATEMENT_ROLLBACK_TRANSACTION = "ROLLBACK TRANSACTION;";

const char *const SQL_STATEMENT_VACUUM = "VACUUM;";

const char *const SQL_STATEMENT_SELECT_SQLITE_MASTER = "SELECT name FROM sqlite_master WHERE type = ? AND name = ?;";
