# EasyBackup

EasyBackup is a backup tool that currently supports xtrabackup.


## Install
```bash
$ make && make install
$ easybackup --help          
Usage:
  easybackup [flags]
  easybackup [command]

Available Commands:
  backup      Take a backup
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Init a repository
  list        List backup sets in repository
  restore     Restore a database from backupset

Flags:
  -h, --help   help for easybackup

Use "easybackup [command] --help" for more information about a command.
```

## Usage

### Backup
1、Create a json file for backup
```bash
$ cat << EOF > config.json
{
    "identifer": "instanceName",
    "version": "8.0.28",
    "login_path": "MYDB8028",
    "db_hostname": "127.0.0.1",
    "db_user": "mysql",
    "throttle": 400, 
    "try_compress": true,
    "bin_path": "/usr/local/xtrabackup/8.0.28/bin",
    "data_path": "/data/mysql/8.0.28",
    "backup_user": "backupuser",
    "backup_hostname": "127.0.0.1"
}
EOF
```

2、Init a repository in path `/data/backup`, the repository name is `repo1`
```bash
$ easybackup init -f config.json -p /data/backup -n repo1
```

3、Take a full backup and check it
```bash
$ easybackup backup -p /data/backup/repo1 -t full
```

```bash
$ easybackup list backupset -p /data/backup/repo1
BackupTime          | Id                                   | Type | FromLSN | ToLSN    | Size(Kb)
2024-01-06 22:58:07 | 25070168-099e-4f6d-9274-bad6edf71ac4 | full | 0       | 18166699 | 3505
```

4、Take a incr backup and check it
```bash
$ easybackup backup -p /data/backup/repo1 -t incr
```

```bash
$ easybackup list backupset -p /data/backup/repo1
BackupTime          | Id                                   | Type | FromLSN  | ToLSN    | Size(Kb)
2024-01-06 22:58:07 | 25070168-099e-4f6d-9274-bad6edf71ac4 | full | 0        | 18166699 | 3505
2024-01-06 23:00:58 | 3d76a621-e9bf-4a67-b663-143923e392ea | incr | 18166699 | 18166719 | 212
```

### Restore
1、Restore database from backupset, the target path is `/data/restore/instance01`
```bash
$ easybackup restore -p /data/backup/repo1 -m /usr/local/mysql/8.0.28 -t /data/restore/instance01 -i 3d76a621-e9bf-4a67-b663-143923e392ea
```

2、Check it, the database is started and port is `36627` !
```bash
$ ps -ef|grep mysql |grep /data/restore/instance01
root      87748      1  0 23:30 ?        00:00:00 sudo -u mysql /usr/local/mysql/8.0.28/bin/mysqld_safe --defaults-file=/data/restore/instance01/my.cnf
mysql     87750  87748  0 23:30 ?        00:00:00 /bin/sh /usr/local/mysql/8.0.28/bin/mysqld_safe --defaults-file=/data/restore/instance01/my.cnf
mysql     87976  87750  0 23:30 ?        00:00:01 /usr/local/mysql/8.0.28/bin/mysqld --defaults-file=/data/restore/instance01/my.cnf --basedir=/usr/local/mysql/8.0.28 --datadir=/data/restore/instance01/25070168-099e-4f6d-9274-bad6edf71ac4 --plugin-dir=/usr/local/mysql/8.0.28/lib/plugin --log-error=/data/restore/instance01/mysql.err --pid-file=/data/restore/instance01/mysql.pid --socket=/data/restore/instance01/mysql.sock --port=36627
```
