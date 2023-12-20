ALTER TABLE public.routes ADD travel_time int NOT NULL DEFAULT 0;
COMMENT ON COLUMN public.routes.travel_time IS 'Время в пути (в минутах)';

ALTER TABLE public.transports ADD registration_number text NOT NULL DEFAULT '';
COMMENT ON COLUMN public.transports.registration_number IS 'Регистрационный номер';
