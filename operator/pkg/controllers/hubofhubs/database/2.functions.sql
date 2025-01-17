CREATE OR REPLACE FUNCTION public.move_applications_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.applications SELECT * FROM spec.applications
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.applications
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;


CREATE OR REPLACE FUNCTION public.move_channels_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.channels SELECT * FROM spec.channels
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.channels
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.move_configs_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.configs SELECT * FROM spec.configs
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.configs
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.move_managedclustersetbindings_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.managedclustersetbindings SELECT * FROM spec.managedclustersetbindings
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.managedclustersetbindings
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.move_managedclustersets_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.managedclustersets SELECT * FROM spec.managedclustersets
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.managedclustersets
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.move_placementbindings_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.placementbindings SELECT * FROM spec.placementbindings
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.placementbindings
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.move_placementrules_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.placementrules SELECT * FROM spec.placementrules
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.placementrules
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.move_placements_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.placements SELECT * FROM spec.placements
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.placements
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.move_policies_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.policies SELECT * FROM spec.policies
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.policies
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.move_subscriptions_to_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  INSERT INTO history.subscriptions SELECT * FROM spec.subscriptions
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  DELETE FROM spec.subscriptions
  WHERE payload -> 'metadata' ->> 'name' = NEW.payload -> 'metadata' ->> 'name' AND
  (
    (
      (payload -> 'metadata' ->> 'namespace' IS NOT NULL AND NEW.payload -> 'metadata' ->> 'namespace' IS NOT NULL)
    AND payload -> 'metadata' ->> 'namespace' = NEW.payload -> 'metadata' ->> 'namespace'
    ) OR (
      payload -> 'metadata' -> 'namespace' IS NULL AND NEW.payload -> 'metadata' -> 'namespace' IS NULL
    )
  );
  RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION public.trigger_set_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;