ALTER TABLE public.seals_data DROP COLUMN rssi;
ALTER TABLE public.seals_data RENAME COLUMN rssi_modem TO rssi;
