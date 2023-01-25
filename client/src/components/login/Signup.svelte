<!-- Component imports and setup -->
<script lang="ts">
  import Surreal from 'surrealdb.js';
  import { loginStore } from '../../stores/login';

  export let toggleSignUp: () => void;

  let email = '';
  let password = '';
  let passwordConfirm = '';
  let error = '';

  // Create a store for the password validation requirements
  const passwordRequirements = {
    length: {
      satisfied: false,
      message: 'At least 5 characters'
    },
    // upperCase: {
    //   satisfied: false,
    //   message: 'At least 1 uppercase letter'
    // },
    lowerCase: {
      satisfied: false,
      message: 'At least 1 lowercase letter'
    },
    // number: {
    //   satisfied: false,
    //   message: 'At least 1 number'
    // }
  };

  // Function to validate the password
  function validatePassword() {
    // Check password length
    if (password.length >= 5) {
      passwordRequirements.length.satisfied = true;
    } else {
      passwordRequirements.length.satisfied = false;
    }

    // // Check for uppercase letter
    // if (/[A-Z]/.test(password)) {
    //   passwordRequirements.upperCase.satisfied = true;
    // } else {
    //   passwordRequirements.upperCase.satisfied = false;
    // }

    // Check for lowercase letter
    if (/[a-z]/.test(password)) {
      passwordRequirements.lowerCase.satisfied = true;
    } else {
      passwordRequirements.lowerCase.satisfied = false;
    }

    // // Check for number
    // if (/\d/.test(password)) {
    //   passwordRequirements.number.satisfied = true;
    // } else {
    //   passwordRequirements.number.satisfied = false;
    // }
  }

  async function handleSignUp() {
    if (signupDisabled) return;
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
      error = err;
      setTimeout(() => error = '', 3000);
    }
  }

  let signupDisabled = true;
  $: {
    signupDisabled = email.length < 1 || password.length < 5 || password !== passwordConfirm;
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
        on:input={validatePassword}
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
        bind:value={passwordConfirm}
        required
      >
    </div>

    <div class="mb-4">
      <div class="text-red-500 font-medium">
        {#each Object.entries(passwordRequirements) as [key, requirement]}
          <div class={requirement.satisfied ? 'text-green-500' : 'text-red-500'}>
            {requirement.satisfied ? '✔' : '✖'} {requirement.message}
          </div>
        {/each}
      </div>
    </div>

    {#if error}
      <div class="text-red-500 text-sm mb-4">{error}</div>
    {/if}
  
    <button
      class={(signupDisabled ? "disabled" : "bg-blue-500 hover:bg-blue-600") + "text-white py-2 px-4 rounded-lg w-full mb-4"}
      disabled={signupDisabled}
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
