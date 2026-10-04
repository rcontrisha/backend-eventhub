ALTER TABLE ONLY public.saved_events DROP CONSTRAINT saved_events_user_id_fkey;
ALTER TABLE ONLY public.saved_events DROP CONSTRAINT saved_events_event_id_fkey;

DROP TABLE public.saved_events;