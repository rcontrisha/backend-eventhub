CREATE TABLE public.communities (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying NOT NULL,
    description text NOT NULL,
    banner_url character varying(255) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

ALTER TABLE ONLY public.communities
    ADD CONSTRAINT communities_pkey PRIMARY KEY (id);