import { render, screen } from '@testing-library/react';
import { SignInWithGoogle } from './SignInWithGoogle';

describe('<SignInWithGoogle />', () => {
    it('renders without crashing', () => {
        const { container } = render(<SignInWithGoogle />);
        expect(container.firstChild).toBeDefined();
    });

    it('displays the correct text', () => {
        render(<SignInWithGoogle />);
        const signInText = screen.getByText(/Sign in with Google/i);
        expect(signInText).toBeInTheDocument();
    });
});