export function ErrorState({ message }: { message: string }) {
  return (
    <div className="flex flex-col items-center justify-center rounded-lg border border-hairline bg-white py-16 text-center">
      <h3 className="font-display text-lg font-semibold text-ink">Something went wrong</h3>
      <p className="mt-1 max-w-sm text-sm text-slate">{message}</p>
    </div>
  );
}