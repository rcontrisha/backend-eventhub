CREATE TABLE public.events (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying(255) NOT NULL,
    "desc" text NOT NULL,
    image_url character varying(255) NOT NULL,
    community_id uuid,
    start_time timestamp without time zone NOT NULL,
    end_time timestamp without time zone NOT NULL,
    location character varying(255) NOT NULL,
    capacity integer,
    organizer_id uuid NOT NULL,
    speakers jsonb DEFAULT '[]'::jsonb
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
);

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_community_id_fkey FOREIGN KEY (community_id) REFERENCES public.communities(id) DEFERRABLE;

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_organizer_id_fkey FOREIGN KEY (organizer_id) REFERENCES public.users(id) DEFERRABLE;