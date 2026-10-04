CREATE TABLE public.testimonials (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    message text NOT NULL,
    job character varying(255) NOT NULL,
    created_at timestamp with time zone NOT NULL
);

ALTER TABLE ONLY public.testimonials
    ADD CONSTRAINT testimonials_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.testimonials
    ADD CONSTRAINT testimonials_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) DEFERRABLE;