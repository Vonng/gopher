-- 通知触发器
CREATE OR REPLACE FUNCTION notify_change() RETURNS TRIGGER AS $$
BEGIN
  IF    (TG_OP = 'INSERT') THEN PERFORM pg_notify(TG_RELNAME || '_chan', 'I' || NEW.id); RETURN NEW;
  ELSIF (TG_OP = 'UPDATE') THEN PERFORM pg_notify(TG_RELNAME || '_chan', 'U' || NEW.id); RETURN NEW;
  ELSIF (TG_OP = 'DELETE') THEN PERFORM pg_notify(TG_RELNAME || '_chan', 'D' || OLD.id); RETURN OLD;
  END IF;
END; $$ LANGUAGE plpgsql SECURITY DEFINER;

-- 用户表
CREATE TABLE users (
  id   TEXT,
  name TEXT,
  PRIMARY KEY (id)
);

--
CREATE TRIGGER t_user_notify AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW EXECUTE PROCEDURE notify_change();