CREATE TABLE public.saved_events (
    event_id uuid NOT NULL,
    user_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);

ALTER TABLE ONLY public.saved_events
    ADD CONSTRAINT saved_events_event_id_fkey FOREIGN KEY (event_id) REFERENCES public.events(id) DEFERRABLE;

ALTER TABLE ONLY public.saved_events
    ADD CONSTRAINT saved_events_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) DEFERRABLE;