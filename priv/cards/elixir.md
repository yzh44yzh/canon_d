# Elixir

## Behaviour

Behaviours provide a way to:
- define a set of functions that have to be implemented by a module;
- ensure that a module implements all the functions in that set.

```
defmodule Parser do
  @callback parse(String.t) :: {:ok, term} | {:error, String.t}
end

defmodule JSONParser do
@behaviour Parser

@impl Parser
  def parse(str), do: {:ok, "some json " <> str}
end
```

## Ecto Associations

One post has many comments.

```
defmodule Post do
  use Ecto.Schema

  schema "posts" do
    field :title, :string
    field :body, :string
    has_many :comments, Comment
  end
end

defmodule Comment do
  use Ecto.Schema

  schema "comments" do
    field :body, :string
    belongs_to :post, Post
  end
end
```

## Ecto Migration: create table

```
defmodule MyApp.Repo.Migrations.CreatePeople do
  use Ecto.Migration

  def change do
    create table(:people) do
      add :first_name, :string
      add :last_name, :string
      add :age, :integer
    end
  end
end
```

## Ecto Migration: add column to table

```
defmodule Repo.Migrations.AlterMyTable do
  use Ecto.Migration

  def change do
    alter table(:my_table) do
      add :new_column, :string
    end
  end
end
```

## Ecto Migration: change column type

```
defmodule Repo.Migrations.AlterMyTable do
  use Ecto.Migration

  def change do
    alter table(:my_table) do
      modify :my_column, :string, from: :text
    end
  end
end
```

## Ecto Migration: rename column

```
defmodule Repo.Migrations.AlterMyTable do
  use Ecto.Migration

  def change do
    rename table(:my_table), :column_old_name, to: :column_new_name
  end
end
```
