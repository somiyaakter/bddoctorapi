export const metadata = {
  title: "API Documentation",
  description: "REST API reference for the MediDirectory doctor directory.",
};

const BASE_URL = "https://bddoctorapi-production.up.railway.app";

export default function ApiDocsPage() {
  return (
    <div className="mx-auto max-w-3xl px-4 py-12 sm:px-6">
      <h1 className="font-display text-3xl font-semibold text-ink">API Documentation</h1>
      <p className="mt-2 text-slate">
        REST API for accessing verified doctor data across Bangladesh. All endpoints
        return JSON and require an API key.
      </p>

      <Section title="Base URL">
        <CodeBlock>{BASE_URL}</CodeBlock>
      </Section>

      <Section title="Authentication">
        <p className="text-sm text-slate">
          Every request must include your API key in the <Code>X-API-Key</Code> header.
        </p>
        <CodeBlock>{`curl "${BASE_URL}/api/v1/doctors" \\\n  -H "X-API-Key: your_api_key"`}</CodeBlock>
        <p className="mt-3 text-sm text-slate">
          Missing or invalid keys return <Code>401 Unauthorized</Code>:
        </p>
        <CodeBlock>{`{ "error": "missing X-API-Key header" }`}</CodeBlock>
      </Section>

      <Section title="Rate limits & quota">
        <ul className="list-disc space-y-2 pl-5 text-sm text-slate">
          <li>Each key has a per-minute request limit. Exceeding it returns <Code>429</Code>.</li>
          <li>
            Each key also has a monthly request quota, reset automatically at the start
            of each calendar month. Exceeding it returns <Code>429</Code> with details:
          </li>
        </ul>
        <CodeBlock>{`{
  "error": "monthly quota exceeded",
  "monthly_quota": 1000,
  "used": 1000,
  "period_start": "2026-08-01",
  "resets_at": "2026-09-01"
}`}</CodeBlock>
      </Section>

      <Section title="GET /api/v1/doctors">
        <p className="text-sm text-slate">Returns a paginated list of doctors.</p>
        <ParamsTable
          rows={[
            ["page", "integer", "Page number. Default 1."],
            ["page_size", "integer", "Results per page. Default 20, max 100."],
            ["q", "string", "Search by name, degrees, specialty, workplace, or BMDC number."],
            ["location_id", "integer", "Filter by location ID (see /api/v1/locations)."],
            ["specialty_id", "integer", "Filter by specialty ID (see /api/v1/specialties)."],
          ]}
        />
        <CodeBlock>{`curl "${BASE_URL}/api/v1/doctors?q=cardio&page=1&page_size=2" \\\n  -H "X-API-Key: your_api_key"`}</CodeBlock>
        <CodeBlock>{`{
  "data": [
    {
      "id": 2,
      "name": "Dr. Md. Shahidullah",
      "bmdc_reg_no": "A-43058",
      "degrees": "MBBS (CMC), BCS (Health), D-Card, MD (BSMMU)",
      "experience_years": 17,
      "specialties": "Cardiology (Heart Diseases) & Medicine Specialist",
      "designation": "Consultant (Cardiology)",
      "workplace": "Chittagong Medical College & Hospital",
      "image_url": "https://www.doctorbangladesh.com/wp-content/uploads/...",
      "profile_url": "https://www.doctorbangladesh.com/dr-mdshahidullah/",
      "chambers": [
        {
          "id": 15073,
          "doctor_id": 2,
          "name": "National Hospital, Chittagong",
          "address": "14/15, Dampara Lane, Mehedibag, Chattogram",
          "visiting_hour": "2.30pm to 5.30pm (Sat, Sun, Tue & Thu)",
          "appointment_phone": "+8801712564193"
        }
      ],
      "created_at": "2026-08-25T07:27:27.700964Z",
      "updated_at": "2026-08-27T14:58:24.313922Z"
    }
  ],
  "pagination": { "page": 1, "page_size": 2, "total_items": 530, "total_pages": 265 }
}`}</CodeBlock>
      </Section>

      <Section title="GET /api/v1/doctors/{id}">
        <p className="text-sm text-slate">
          Returns a single doctor by ID. Returns <Code>404</Code> with{" "}
          <Code>{`{ "error": "doctor not found" }`}</Code> 
        </p>
        <CodeBlock>{`curl "${BASE_URL}/api/v1/doctors/2" \\\n  -H "X-API-Key: your_api_key"`}</CodeBlock>
      </Section>

      <Section title="GET /api/v1/locations">
        <p className="text-sm text-slate">Returns all locations (divisions/cities).</p>
        <CodeBlock>{`curl "${BASE_URL}/api/v1/locations" \\\n  -H "X-API-Key: your_api_key"`}</CodeBlock>
        <CodeBlock>{`{ "data": [{ "id": 1, "name": "Dhaka", "url": "https://www.doctorbangladesh.com/doctors-dhaka/" }] }`}</CodeBlock>
      </Section>

      <Section title="GET /api/v1/specialties">
        <p className="text-sm text-slate">
          Returns all specialties, or only those in a given location.
        </p>
        <ParamsTable
          rows={[["location_id", "integer", "Optional. Filter specialties by location."]]}
        />
        <CodeBlock>{`curl "${BASE_URL}/api/v1/specialties?location_id=1" \\\n  -H "X-API-Key: your_api_key"`}</CodeBlock>
      </Section>
    </div>
  );
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <section className="mt-10 border-t border-hairline pt-8">
      <h2 className="font-display text-lg font-semibold text-ink">{title}</h2>
      <div className="mt-3 space-y-3">{children}</div>
    </section>
  );
}

function Code({ children }: { children: React.ReactNode }) {
  return (
    <code className="rounded bg-secondary px-1.5 py-0.5 font-mono text-xs text-ink">
      {children}
    </code>
  );
}

function CodeBlock({ children }: { children: string }) {
  return (
    <pre className="overflow-x-auto rounded-md border border-hairline bg-ink px-4 py-3 text-xs text-paper">
      <code>{children}</code>
    </pre>
  );
}

function ParamsTable({ rows }: { rows: [string, string, string][] }) {
  return (
    <div className="overflow-x-auto rounded-md border border-hairline">
      <table className="w-full text-left text-sm">
        <thead className="bg-secondary text-xs uppercase text-slate">
          <tr>
            <th className="px-3 py-2">Param</th>
            <th className="px-3 py-2">Type</th>
            <th className="px-3 py-2">Description</th>
          </tr>
        </thead>
        <tbody>
          {rows.map(([param, type, desc]) => (
            <tr key={param} className="border-t border-hairline">
              <td className="px-3 py-2">
                <Code>{param}</Code>
              </td>
              <td className="px-3 py-2 text-slate">{type}</td>
              <td className="px-3 py-2 text-slate">{desc}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}