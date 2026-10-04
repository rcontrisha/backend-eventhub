CREATE TABLE public.event_discussions (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    event_id uuid NOT NULL,
    user_id uuid NOT NULL,
    message text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

ALTER TABLE ONLY public.event_discussions
    ADD CONSTRAINT event_discussions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.event_discussions
    ADD CONSTRAINT event_discussions_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) DEFERRABLE;

ALTER TABLE ONLY public.event_discussions
    ADD CONSTRAINT event_discussions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) DEFERRABLE;