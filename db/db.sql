
CREATE SCHEMA users

CREATE TABLE users.users(
user_id int not null,
is_active boolean not null,
first_name varchar(50),
last_name varchar(50),
email varchar(50) not null,
pswd_hash_bytes bytea not null
)

ALTER TABLE users.users ADD CONSTRAINT PK_users PRIMARY KEY(user_id)

CREATE INDEX ON users.users(email)

CREATE OR REPLACE FUNCTION users.user_getByEmail(
	_email varchar(50))
    RETURNS TABLE(user_id integer, is_active boolean, first_name character varying, last_name character varying, email character varying) 
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE 
    ROWS 1000
AS $BODY$	
BEGIN

RETURN QUERY
	SELECT	u.user_id
			, u.is_active
			, u.first_name
			, u.last_name
			, u.email
	FROM	users.users u
	WHERE	u.email = _email;
			
END;
$BODY$;

CREATE OR REPLACE FUNCTION users.user_get(
	_user_id integer)
    RETURNS TABLE(user_id integer, is_active boolean, first_name character varying, last_name character varying, email character varying) 
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE 
    ROWS 1000
AS $BODY$	
BEGIN

RETURN QUERY
	SELECT	u.user_id
			, u.is_active
			, u.first_name
			, u.last_name
			, u.email
	FROM	users.users u
	WHERE	u.user_id = _user_id;
			
END;
$BODY$;

CREATE OR REPLACE FUNCTION users.user_getbyauth(
	_email character varying,
	_pswd_hash_bytes bytea)
    RETURNS TABLE(user_id integer, is_active boolean, first_name character varying, last_name character varying, email character varying) 
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE 
    ROWS 1000
AS $BODY$	
BEGIN

RETURN QUERY
	SELECT	u.user_id
			, u.is_active
			, u.first_name
			, u.last_name
			, u.email
	FROM	users.users u
	WHERE	u.email = _email
			AND u.pswd_hash_bytes = _pswd_hash_bytes;
			
END;
$BODY$;

CREATE SEQUENCE users.user_id
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

CREATE OR REPLACE FUNCTION users.user_save(
	_user_id integer,
	_is_active boolean,
	_first_name character varying,
	_last_name character varying,
	_email character varying,
	_pswd_hash_bytes bytea)
    RETURNS TABLE(user_id integer) 
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE 
    ROWS 1000
AS $BODY$
	DECLARE _user_id_tmp int;
	DECLARE _error smallint;
BEGIN

	_user_id_tmp := _user_id;

	IF _user_id_tmp = 0 THEN
		_user_id_tmp := nextval('users.user_id');
	END IF;
	
	IF	EXISTS	(
			SELECT	1
			FROM	users.users u
			WHERE	u.email = _email
					AND u.user_id != _user_id_tmp
		)
	THEN			
		RAISE EXCEPTION 'E-Mail уже используется! %', _error USING ERRCODE = '22024'; 
	END IF;
	
	INSERT INTO users.users( user_id
							, is_active
							, first_name
							, last_name
							, email
						    , pswd_hash_bytes)
	SELECT	_user_id_tmp
			, _is_active
			, _first_name
			, _last_name
			, _email
			, _pswd_hash_bytes
	ON CONFLICT ON CONSTRAINT PK_users 
	DO UPDATE SET	is_active = _is_active
					, first_name = _first_name
					, last_name = _last_name
					, email = _email
					, pswd_hash_bytes = _pswd_hash_bytes;
					
RETURN QUERY 
	SELECT	_user_id_tmp user_id;
		
END;
$BODY$;