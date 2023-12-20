ALTER TABLE public.modems_log_raw ADD remote_port int NOT NULL DEFAULT 0;
COMMENT ON COLUMN public.modems_log_raw.remote_port IS 'Порт устройства';