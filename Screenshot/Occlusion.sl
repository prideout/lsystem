class Occlusion(float darkness = 1;
                color em = (1,0,1);
                float samples = 128;
                float maxdist = 100;
)
{
  public void surface(output color Ci, Oi)
  {
    normal Nf = normalize(faceforward(N, I));
    normal Nn = normalize(Nf);
    float occ = darkness * occlusion(P, Nn, samples, "maxdist", maxdist);
    Ci = (1 - occ) * Cs * em * Os;
    Oi = Os;
  }
}
