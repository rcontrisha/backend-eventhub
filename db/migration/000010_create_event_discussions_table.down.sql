ALTER TABLE ONLY public.event_discussions DROP CONSTRAINT event_discussions_user_id_fkey;
ALTER TABLE ONLY public.event_discussions DROP CONSTRAINT event_discussions_event_id_fkey;
ALTER TABLE ONLY public.event_discussions DROP CONSTRAINT event_discussions_pkey;

DROP TABLE public.event_discussions;