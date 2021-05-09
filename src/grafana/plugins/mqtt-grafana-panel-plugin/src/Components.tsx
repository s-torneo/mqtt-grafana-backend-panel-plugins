import React from 'react';

function Button(props: any) {
  return (
    <button
      className={props.classname}
      style={{ color: props.textcolor, background: props.backgroundcolor }}
      onClick={props.handle}
    >
      {props.title}
    </button>
  );
}

function Form(props: any) {
  return (
    <label>
      {props.name}
      <input
        className="form"
        type={props.type}
        value={props.value}
        placeholder={props.placeholder}
        onChange={(e) => {
          props.handle(e);
        }}
      />
    </label>
  );
}

export { Button, Form };
