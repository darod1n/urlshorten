package postgresql

import (
	"context"
	"database/sql"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	urls map[string]string
	path string
	mu   *sync.Mutex
	base *sql.DB
}

func (db *DB) AddURL(ctx context.Context, url string) (string, error) {
	row := db.base.QueryRowContext(ctx, "INSERT INTO urls (original_url) VALUES($1) returning short_url;", url)
	var shortURL string
	row.Scan(&shortURL)
	return shortURL, nil
}

func (db *DB) GetURL(ctx context.Context, shortURL string) (string, error) {
	row := db.base.QueryRowContext(ctx, "select original_url from urls where short_url=$1;", shortURL)
	var originalURL string
	if err := row.Scan(&originalURL); err != nil {
		return "", err
	}
	return originalURL, nil
}

func (db *DB) PingContext(ctx context.Context) error {
	return db.base.PingContext(ctx)
}

func (db *DB) Close() error {
	return db.base.Close()
}

func createDB(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
	create table if not exists urls(
		short_url text primary key,
		original_url text
	);
	
	Create or replace function random_string() returns text as
	$$
	declare
	  chars text[] := '{0,1,2,3,4,5,6,7,8,9,A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z}';
	  result text := '';
	  i integer := 0;
	  length integer := 6;
	begin
	  if length < 0 then
		raise exception 'Given length cannot be less than 0';
	  end if;
	  for i in 1..length loop
		result := result || chars[1+random()*(array_length(chars, 1)-1)];
	  end loop;
	  return result;
	end;
	$$ language plpgsql;
	
	CREATE OR REPLACE FUNCTION unique_short_id()
	RETURNS TRIGGER AS $$
	
	 -- Declare the variables we'll be using.
	DECLARE
	  key TEXT;
	  qry TEXT;
	  found TEXT;
	BEGIN
	
	  qry := 'SELECT short_url FROM ' || quote_ident(TG_TABLE_NAME) || ' WHERE short_url=';
	  LOOP
		key := random_string();
	
		EXECUTE qry || quote_literal(key) INTO found;
		IF found IS NULL THEN
		  EXIT;
		END IF;
	
	  END LOOP;
	
	  
	  NEW.short_url = key;
	  RETURN NEW;
	END;
	$$ language 'plpgsql';
	
	CREATE OR REPLACE TRIGGER trigger_urls_genid BEFORE INSERT ON urls FOR EACH ROW EXECUTE PROCEDURE unique_short_id();
	`)
	if err != nil {
		return err
	}
	return nil
}

func NewDB(path, driverName, dataSourceName string) (*DB, error) {

	base, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if err := createDB(ctx, base); err != nil {
		return nil, err
	}

	return &DB{
		mu:   &sync.Mutex{},
		path: path,
		base: base,
	}, nil
}