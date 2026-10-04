ALTER TABLE ONLY public.event_tags DROP CONSTRAINT event_tags_tag_id_fkey;
ALTER TABLE ONLY public.event_tags DROP CONSTRAINT event_tags_event_id_fkey;

DROP TABLE public.event_tags;