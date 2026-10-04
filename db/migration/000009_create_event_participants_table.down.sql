ALTER TABLE ONLY public.event_participants DROP CONSTRAINT event_participants_user_id_fkey;
ALTER TABLE ONLY public.event_participants DROP CONSTRAINT event_participants_event_id_fkey;

DROP TABLE public.event_participants;