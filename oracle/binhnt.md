# Config file 
 /etc/sysconfig/oracledb_ORCLCDB-19c.conf 


# Init connection 
export ORACLE_HOME=/opt/oracle/product/19c/dbhome_1
export PATH=$ORACLE_HOME/bin:$PATH
export ORACLE_SID=ORCLCDB
sqlplus / as sysdba

# Create user 
sqlplus as sysdba <<EOF
alter session set container=ORCLPDB1;
create user testuser1 identified by testuser1 quota unlimited on users;
grant connect, resource to testuser1;
exit;
EOF

# Access db 

SELECT * FROM dba_users;
select owner FROM all_tables group by owner;
select distinct owner FROM all_tables;

SELECT * FROM all_tables WHERE OWNER  LIKE '%SCHEMA_NAME%';
SELECT * FROM all_tables WHERE TABLE_NAME  LIKE '%TABLE_NAME%'


select username as schema_name
from sys.all_users
order by username;

select username as schema_name
from sys.dba_users
order by username;

SELECT TABLESPACE_NAME FROM USER_TABLESPACES;
SELECT USERNAME FROM ALL_USERS ORDER BY USERNAME; 


SELECT PDB_ID, PDB_NAME, STATUS FROM DBA_PDBS ORDER BY PDB_ID;

SELECT owner, table_name FROM dba_tables ;

SELECT * FROM all_objects WHERE object_type IN (‘TABLE’,’VIEW’) AND object_name = ‘OBJECT_NAME’;


# Access db 
sqlplus testuser1/testuser1@//localhost:1521/ORCLPDB1

sqlplus testuser1/testuser1@//localhost:1521/ORCLPDB1



# Add lib 
export LD_LIBRARY_PATH=$HOME/Downloads/instantclient_19_8:$LD_LIBRARY_PATH



