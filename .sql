CREATE TABLE IF NOT EXISTS public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    create_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    login TEXT NOT NULL,
    password_hash TEXT NOT NULL, 
    balance TEXT NOT NULL,
    withdrawn_amount TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS public.orders (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    status text NOT NULL DEFAULT 'NEW',
    accrual TEXT NOT NULL,
    sum TEXT NOT NULL,
    upload_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	constraint user_id foreign key (user_id) references public.users (id)
);

CREATE TABLE IF NOT EXISTS public.withdrawals (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    sum TEXT NOT NULL,
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	constraint user_id foreign key (user_id) references public.users (id)
);

