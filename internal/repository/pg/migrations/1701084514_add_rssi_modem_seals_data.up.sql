ALTER TABLE public.seals_data ADD rssi_modem int2 NOT NULL DEFAULT 0;
COMMENT ON COLUMN public.seals_data.rssi_modem IS 'RSSI базового блока, дБм';