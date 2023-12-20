ALTER TABLE public.modems_data ADD signal_gps int2 NOT NULL DEFAULT 0;
ALTER TABLE public.modems_data ADD signal_glonass int2 NOT NULL DEFAULT 0;

ALTER TABLE public.coordinates ADD signal_gps int2 NOT NULL DEFAULT 0;
ALTER TABLE public.coordinates ADD signal_glonass int2 NOT NULL DEFAULT 0;