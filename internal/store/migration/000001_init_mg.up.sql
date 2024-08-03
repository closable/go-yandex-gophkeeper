CREATE SCHEMA IF NOT EXISTS gophkeeper;

CREATE TABLE IF NOT EXISTS gophkeeper.users (
				user_id bigserial NOT NULL,
				login varchar(250) NOT NULL,
				password varchar(250) NOT NULL,
				is_active bool DEFAULT true NULL,
				key varchar(255) NULL,
				CONSTRAINT users_pkey PRIMARY KEY (user_id)
				);
				
CREATE TABLE IF NOT EXISTS gophkeeper.users_data (
				id bigserial NOT NULL,
				user_id bigserial NOT NULL,
				data_type int4 NULL,
				data text NULL,
				is_deleted bool DEFAULT false NULL,
				name varchar(255) NULL,
				is_restore bool DEFAULT false NULL,
				CONSTRAINT users_data_pkey PRIMARY KEY (id)
			);

CREATE TABLE IF NOT EXISTS gophkeeper.data_types (
				id bigserial NOT NULL,
				name varchar(255) NOT NULL,
				is_deleted bool DEFAULT false NULL,
				CONSTRAINT data_types_pkey PRIMARY KEY (id)
			);
            
merge into gophkeeper.data_types dt	using (
					select 1 id, 'PlainText' name,  false is_deleted 
					union
					select 2 id, 'KeyValue' name, false is_deleted
					union 
					select 3 id, 'FileData' name, false is_deleted
					union 
					select 4 id, 'FolderData' name, false is_deleted
				) as res on (dt.id = res.id)
					when not matched then 
					insert (id, name, is_deleted)
					values (res.id, res.name, res.is_deleted);