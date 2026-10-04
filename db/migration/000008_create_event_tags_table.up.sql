CREATE TABLE public.event_tags (
    event_id uuid NOT NULL,
    tag_id uuid NOT NULL
);

ALTER TABLE ONLY public.event_tags
    ADD CONSTRAINT event_tags_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) DEFERRABLE;

ALTER TABLE ONLY public.event_tags
    ADD CONSTRAINT event_tags_tag_id_fkey FOREIGN KEY (tag_id) REFERENCES public.tags(id) DEFERRABLE;