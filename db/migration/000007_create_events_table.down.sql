ALTER TABLE ONLY public.events DROP CONSTRAINT events_organizer_id_fkey;
ALTER TABLE ONLY public.events DROP CONSTRAINT events_community_id_fkey;
ALTER TABLE ONLY public.events DROP CONSTRAINT events_pkey;

DROP TABLE public.events;