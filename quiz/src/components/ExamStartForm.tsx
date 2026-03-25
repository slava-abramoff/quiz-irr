type FormFields = {
  fullName: string;
  email: string;
  org: string;
  birthYear: string;
};

type ExamStartFormProps = {
  values: FormFields;
  onChange: (field: keyof FormFields, value: string) => void;
  onSubmit: () => void;
  loading: boolean;
  error: string | null;
};

export function ExamStartForm({
  values,
  onChange,
  onSubmit,
  loading,
  error,
}: ExamStartFormProps) {
  const filled =
    values.fullName.trim().length > 0 &&
    values.email.trim().length > 0 &&
    values.org.trim().length > 0 &&
    values.birthYear.trim().length > 0;

  return (
    <form
      className="exam-form"
      onSubmit={(e) => {
        e.preventDefault();
        if (filled && !loading) onSubmit();
      }}
    >
      <h2 className="exam-form-title">Данные для начала</h2>
      <label className="exam-field">
        <span className="exam-field-label">ФИО</span>
        <input
          className="exam-input"
          type="text"
          name="fullName"
          autoComplete="name"
          value={values.fullName}
          onChange={(e) => onChange("fullName", e.target.value)}
          disabled={loading}
        />
      </label>
      <label className="exam-field">
        <span className="exam-field-label">Ваш возраст</span>
        <select
          className="exam-input"
          name="birthYear"
          value={values.birthYear}
          onChange={(e) => onChange("birthYear", e.target.value)}
          disabled={loading}
        >
          <option value="" disabled>
            Выберите год рождения
          </option>
          <option value="2008">2008</option>
          <option value="2009">2009</option>
          <option value="2010">2010</option>
          <option value="2011">2011</option>
          <option value="2012">2012</option>
          <option value="2013">2013</option>
        </select>
      </label>
      <label className="exam-field">
        <span className="exam-field-label">Почта</span>
        <input
          className="exam-input"
          type="email"
          name="email"
          autoComplete="email"
          value={values.email}
          onChange={(e) => onChange("email", e.target.value)}
          disabled={loading}
        />
      </label>
      <label className="exam-field">
        <span className="exam-field-label">Учебное учреждение</span>
        <input
          className="exam-input"
          type="text"
          name="org"
          autoComplete="organization"
          value={values.org}
          onChange={(e) => onChange("org", e.target.value)}
          disabled={loading}
        />
      </label>
      {error ? (
        <p className="exam-form-error" role="alert">
          {error}
        </p>
      ) : null}
      <button
        type="submit"
        className="quiz-start-btn"
        disabled={!filled || loading}
      >
        {loading ? "Отправка…" : "Старт"}
      </button>
    </form>
  );
}
