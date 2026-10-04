CREATE TABLE public.event_participants (
    event_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);

ALTER TABLE ONLY public.event_participants
    ADD CONSTRAINT event_participants_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) DEFERRABLE;

ALTER TABLE ONLY public.event_participants
    ADD CONSTRAINT event_participants_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) DEFERRABLE;