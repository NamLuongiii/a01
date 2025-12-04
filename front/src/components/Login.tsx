import { useMutation } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { MUTATION_KEYS } from "../constants";

type LoginForm = {
  name: string;
};

export default function Login() {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>();

  const submit = async (data: LoginForm) => {
    // default behavior for demo
    mutateAsync(data);
  };

  // query login
  const { mutateAsync, isPending } = useMutation({
    mutationKey: [MUTATION_KEYS.LOGIN],
    mutationFn: async (data: LoginForm) => {
      // perform login logic here
      console.log("Performing login with:", data);
    },
    onError: (error) => {
      alert(error.message);
    },
  });

  return (
    <form onSubmit={handleSubmit(submit)} noValidate style={{ maxWidth: 360 }}>
      <div style={{ marginBottom: 12 }}>
        <label htmlFor="name" style={{ display: "block", marginBottom: 4 }}>
          Name
        </label>
        <input
          id="name"
          type="text"
          {...register("name", {
            required: "Name is required",
          })}
          aria-invalid={errors.name ? "true" : "false"}
        />
        {errors.name && (
          <p role="alert" style={{ color: "crimson", marginTop: 6 }}>
            {errors.name.message}
          </p>
        )}
      </div>

      <button
        type="submit"
        disabled={isPending}
        style={{
          padding: "8px 12px",
          cursor: isPending ? "not-allowed" : "pointer",
        }}
      >
        {isPending ? "Signing in..." : "Sign in"}
      </button>
    </form>
  );
}
