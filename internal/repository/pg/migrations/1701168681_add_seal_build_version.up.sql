ALTER TABLE public.seals_data ADD build_version int4 NOT NULL DEFAULT 0;
COMMENT ON COLUMN public.seals_data.build_version IS 'Build версия пломбы';
ALTER TABLE public.modems_data DROP COLUMN records_count;