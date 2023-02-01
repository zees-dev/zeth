<!-- Component imports and setup -->
<script lang="ts">
  import Surreal from 'surrealdb.js';
  import { z } from 'zod';
  import { loginStore } from '../../stores/login';

  export let toggleSignUp: () => void;

  const formSchema = z.object({
    email: z.string({ required_error: "Email is required" })
      .email("Invalid email address")
      .min(3, "Email must be atleast 3 characters")
      .max(64, "Email must be atmost 64 characters")
      .trim(),
    password: z.string({ required_error: "Password is required"})
      .min(5, "Password must be atleast 5 characters")
      .max(64, "Password must be atmost 64 characters")
      // .regex(/[0-9]/, "Password must contain atleast 1 number")
      // .regex(/[A-Z]/, "Password must contain atleast 1 uppercase letter")
      .regex(/[a-z]/, "Password must contain atleast 1 lowercase letter"),
    confirmPassword: z.string(),
    // toc: z.enum(['on']),
  }).superRefine(({ confirmPassword, password }, ctx) => {
    if (confirmPassword !== password) {
      ctx.addIssue({
        code: "custom",
        message: "The passwords did not match",
        path: ['confirmPassword'],
      });
    }
  });

  let email = '';
  let password = '';
  let confirmPassword = '';
  let error = '';

  async function handleSignUp() {
    if (validationErrors.length) return;
    try {
      const token = await Surreal.Instance.signup({
        NS: 'test',
        DB: 'test',
        SC: 'allusers',
        user: email.toLowerCase(),
        pass: password,
        tags: ['zeth'],
      });
      console.info('signed up to db: ', token);
      toggleSignUp(); // for next login
      loginStore.login(token);
    } catch (err) {
      error = err as string;
      setTimeout(() => error = '', 3000);
    }
  }

  let validationErrors: string[] = [];
  $: {
    try {
      formSchema.parse({ email, password, confirmPassword })
      validationErrors = []
    } catch (err) {
      validationErrors = Object.values((err as any).flatten().fieldErrors).flat() as string[]
    }
  }
</script>

<!-- Component template -->
<div class="flex justify-center items-center h-screen">
  <form class="bg-white p-6 rounded-lg flex flex-col" on:submit|preventDefault={handleSignUp}>
    <h2 class="text-lg font-medium mb-4 text-center">Sign Up</h2>
    <div class="mb-4">
      <label class="block text-gray-700 font-medium mb-2" for="email">
        Email
      </label>
      <input
        class="border border-gray-400 p-2 rounded-lg w-full"
        type="email"
        id="email"
        placeholder="Enter your email"
        bind:value={email}
        required
      >
    </div>

    <div class="mb-4">
      <label class="block text-gray-700 font-medium mb-2" for="password">
        Password
      </label>
      <input
        class="border border-gray-400 p-2 rounded-lg w-full"
        type="password"
        id="password"
        bind:value={password}
        required
      >
    </div>

    <div class="mb-4">
      <label class="block text-gray-700 font-medium mb-2" for="password-confirm">
      Confirm Password
      </label>
      <input
        class="border border-gray-400 p-2 rounded-lg w-full"
        type="password"
        id="password-confirm"
        bind:value={confirmPassword}
        required
      >
    </div>

    <div class="mb-4">
      <div class="text-red-500 font-normal">
        {#if !formSchema.safeParse({ email, password, confirmPassword }).success}
          {#each validationErrors as error }
            <div class="text-red-500">
              ❌ {error}
            </div>
          {/each}
        {/if}
      </div>
    </div>

    {#if error}
      <div class="text-red-500 text-sm mb-4">{error}</div>
    {/if}
  
    <button
      class={(validationErrors.length > 0 ? "disabled" : "bg-blue-500 hover:bg-blue-600") + "text-white py-2 px-4 rounded-lg w-full mb-4"}
      disabled={validationErrors.length > 0}
    >
      Sign Up
    </button>
  
    <button class="bg-white text-gray-600 py-2 px-4 rounded-lg hover:bg-gray-200 w-full" on:click={toggleSignUp}>
      Log In
    </button>
  </form>
</div>

<style>
  form {
    width: 20rem;
  }
</style>
